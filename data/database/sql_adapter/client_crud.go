package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CreateClientTable creates the client table in the database
// Returns any errors
func (crud *SQLCRUD) CreateClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create client table script", err)
	}

	return err
}

// DropClientTable drops the client table from the database
// Returns any errors
func (crud *SQLCRUD) DropClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop client table script", err)
	}

	return err
}

// SaveClient validates the client model is valid and inserts a new row into the client table
// Returns any errors
func (crud *SQLCRUD) SaveClient(client *models.Client) error {
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveClientScript(), client.ID, client.Name)
	cancel()

	if err != nil {
		return common.ChainError("error executing save client statement", err)
	}

	return nil
}

// UpdateClient validates the client model is valid and updates the row in the client table
// Returns result of whether the client was found, and any errors
func (crud *SQLCRUD) UpdateClient(client *models.Client) (bool, error) {
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return false, errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateClientScript(), client.ID, client.Name)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing save client statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

// DeleteUser deletes the row in the user table with the matching id
// Returns result of whether the client was found, and any errors
func (crud *SQLCRUD) DeleteClient(id uuid.UUID) (bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteClientScript(), id)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing delete client statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

// GetClientByID gets the row in the client table with the matching id, and creates a new client model using its data
// Returns the model and any errors
func (crud *SQLCRUD) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetClientByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get client by id query", err)
	}
	defer rows.Close()

	return readClientData(rows)
}

func readClientData(rows *sql.Rows) (*models.Client, error) {
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
	client := &models.Client{}
	err := rows.Scan(&client.ID, &client.Name)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return client, nil
}
