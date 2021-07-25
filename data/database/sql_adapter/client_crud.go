package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (crud *SQLCRUD) CreateClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create client table script", err)
	}

	return err
}

func (crud *SQLCRUD) DropClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop client table script", err)
	}

	return err
}

// SaveClient validates the client model is valid and inserts a new row into the client table.
// Returns any errors.
func (crud *SQLCRUD) SaveClient(client *models.Client) error {
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveClientScript(), client.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing save client statement", err)
	}

	return nil
}

// GetClientByID gets the row in the client table with the matching id, and creates a new client model using its data.
// Returns the model and any errors.
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
	err := rows.Scan(&client.ID)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return client, nil
}
