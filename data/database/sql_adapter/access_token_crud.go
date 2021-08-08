package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CreateAccessTokenTable creates the access token table in the database
// Returns any errors
func (crud *SQLCRUD) CreateAccessTokenTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateAccessTokenTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create access token table script", err)
	}

	return err
}

// DropAccessTokenTable drops the access token table from the database
// Returns any errors
func (crud *SQLCRUD) DropAccessTokenTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropAccessTokenTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop access token table script", err)
	}

	return err
}

// SaveAccessToken validates the access token model is valid and inserts a new row into the access_token table
// Returns any errors
func (crud *SQLCRUD) SaveAccessToken(token *models.AccessToken) error {
	verr := token.Validate()
	if verr != models.ValidateAccessTokenValid {
		return errors.New(fmt.Sprint("error validating access token model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveAccessTokenScript(),
		token.ID, token.User.ID, token.Client.ID, token.Scope.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing save access token statement", err)
	}

	return nil
}

// GetAccessTokenByID gets the row in the access_token table with the matching id, and creates a new access token model with associated models using its data
// Returns the model and any errors
func (crud *SQLCRUD) GetAccessTokenByID(ID uuid.UUID) (*models.AccessToken, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetAccessTokenByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get access token by id query", err)
	}
	defer rows.Close()

	return readAccessTokenData(rows)
}

// DeleteAccessToken deletes the row in the access_token table with the matching id
// Returns any errors
func (crud *SQLCRUD) DeleteAccessToken(token *models.AccessToken) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteAccessTokenScript(), token.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete access token statement", err)
	}

	return nil
}

// DeleteAllOtherUserTokens deletes all the rows in the access_token table with the matching user id, and not the token id
// Returns any errors
func (crud *SQLCRUD) DeleteAllOtherUserTokens(token *models.AccessToken) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteAllOtherUserTokensScript(), token.User.ID, token.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete all other user tokens statement", err)
	}

	return nil
}

func readAccessTokenData(rows *sql.Rows) (*models.AccessToken, error) {
	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, common.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	token := &models.AccessToken{
		User:   &models.User{},
		Client: &models.Client{},
		Scope:  &models.Scope{},
	}

	//get the result
	err := rows.Scan(
		&token.ID,
		&token.User.ID, &token.User.Username, &token.User.PasswordHash,
		&token.Client.ID, &token.Client.Name,
		&token.Scope.ID, &token.Scope.Name,
	)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return token, nil
}
