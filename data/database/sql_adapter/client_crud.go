package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CreateClientTable creates the client table in the database.
// Returns any errors.
func (crud *SQLCRUD) CreateClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create client table script", err)
	}

	return err
}

// DropClientTable drops the client table from the database.
// Returns any errors.
func (crud *SQLCRUD) DropClientTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropClientTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop client table script", err)
	}

	return err
}

// CreateClient validates the client model is valid and inserts a new row into the client table.
// Returns any errors.
func (crud *SQLCRUD) CreateClient(client *models.Client) error {
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateClientScript(), client.UID, client.Name, client.RedirectUrl)
	cancel()

	if err != nil {
		return common.ChainError("error executing create client statement", err)
	}

	return nil
}

// GetClientByUID gets the row in the client table with the matching uid, and creates a new client model using its data.
// Returns the model and any errors.
func (crud *SQLCRUD) GetClientByUID(uid uuid.UUID) (*models.Client, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetClientByUIDScript(), uid)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get client by uid query", err)
	}
	defer rows.Close()

	return readClientData(rows)
}

// UpdateClient validates the client model is valid and updates the row in the client table.
// Returns result of whether the client was found, and any errors.
func (crud *SQLCRUD) UpdateClient(client *models.Client) (bool, error) {
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return false, errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateClientScript(), client.UID, client.Name, client.RedirectUrl)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update client statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

// DeleteUser deletes the row in the user table with the matching uid.
// Returns result of whether the client was found, and any errors.
func (crud *SQLCRUD) DeleteClient(uid uuid.UUID) (bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteClientScript(), uid)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing delete client statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
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
	err := rows.Scan(&client.UID, &client.Name, &client.RedirectUrl)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return client, nil
}
