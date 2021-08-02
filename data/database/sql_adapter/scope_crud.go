package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"
)

// CreateScopeTable creates the scope table in the database
// Returns any errors
func (crud *SQLCRUD) CreateScopeTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateScopeTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create scope table script", err)
	}

	return err
}

// DropScopeTable drops the scope table from the database
// Returns any errors
func (crud *SQLCRUD) DropScopeTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropScopeTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop scope table script", err)
	}

	return err
}

// SaveScope validates the scope model is valid and inserts a new row into the scope table
// Returns any errors
func (crud *SQLCRUD) SaveScope(scope *models.Scope) error {
	verr := scope.Validate()
	if verr != models.ValidateScopeValid {
		return errors.New(fmt.Sprint("error validating scope model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveScopeScript(), scope.ID, scope.Name)
	cancel()

	if err != nil {
		return common.ChainError("error executing save scope statement", err)
	}

	return nil
}

// GetScopeByName gets the row in the scope table with the matching name, and creates a new scope model using its data
// Returns the scope and any errors
func (crud *SQLCRUD) GetScopeByName(name string) (*models.Scope, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetScopeByNameScript(), name)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get scope by name query", err)
	}
	defer rows.Close()

	return readScopeData(rows)
}

func readScopeData(rows *sql.Rows) (*models.Scope, error) {
	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, common.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	//get the result
	scope := &models.Scope{}
	err := rows.Scan(&scope.ID, &scope.Name)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return scope, nil
}
