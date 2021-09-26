package sqladapter

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"

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

func (crud *SQLCRUD) CreateClient(client *models.Client) error {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateClientScript(),
		client.UID, client.Name, client.RedirectUrl, client.TokenType, client.KeyUri)
	cancel()

	if err != nil {
		return common.ChainError("error executing create client statement", err)
	}

	return nil
}

func (crud *SQLCRUD) GetClients() ([]*models.Client, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetClientsScript())
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get clients query", err)
	}
	defer rows.Close()

	//read the data
	clients := []*models.Client{}
	for {
		client, err := readClientData(rows)
		if err != nil {
			return nil, err
		}

		if client == nil {
			break
		}
		clients = append(clients, client)
	}
	return clients, nil
}

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

func (crud *SQLCRUD) UpdateClient(client *models.Client) (bool, error) {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return false, errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateClientScript(),
		client.UID, client.Name, client.RedirectUrl, client.TokenType, client.KeyUri)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update client statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

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
	err := rows.Scan(&client.UID, &client.Name, &client.RedirectUrl, &client.TokenType, &client.KeyUri)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return client, nil
}
