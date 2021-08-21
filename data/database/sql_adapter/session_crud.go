package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CreateSessionTable creates the session table in the database.
// Returns any errors.
func (crud *SQLCRUD) CreateSessionTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateSessionTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create session table script", err)
	}

	return err
}

// DropSessionTable drops the session table from the database.
// Returns any errors.
func (crud *SQLCRUD) DropSessionTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropSessionTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop session table script", err)
	}

	return err
}

// SaveSession validates the session model is valid and inserts a new row into the session table.
// Returns any errors.
func (crud *SQLCRUD) SaveSession(token *models.Session) error {
	verr := token.Validate()
	if verr != models.ValidateSessionValid {
		return errors.New(fmt.Sprint("error validating session model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveSessionScript(), token.ID, token.User.Username)
	cancel()

	if err != nil {
		return common.ChainError("error executing save session statement", err)
	}

	return nil
}

// GetSessionByID gets the row in the session table with the matching id, and creates a new session model with associated models using its data.
// Returns the model and any errors.
func (crud *SQLCRUD) GetSessionByID(ID uuid.UUID) (*models.Session, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetSessionByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get session by id query", err)
	}
	defer rows.Close()

	return readSessionData(rows)
}

// DeleteSession deletes the row in the session table with the matching id.
// Returns any errors.
func (crud *SQLCRUD) DeleteSession(token *models.Session) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteSessionScript(), token.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete session statement", err)
	}

	return nil
}

// DeleteAllOtherUserSessions deletes all the rows in the session table with the matching user id, and not the token id.
// Returns any errors.
func (crud *SQLCRUD) DeleteAllOtherUserSessions(token *models.Session) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteAllOtherUserSessionsScript(), token.ID, token.User.Username)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete all other user sessions statement", err)
	}

	return nil
}

func readSessionData(rows *sql.Rows) (*models.Session, error) {
	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, common.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	token := &models.Session{
		User: &models.User{},
	}

	//get the result
	err := rows.Scan(
		&token.ID,
		&token.User.Username, &token.User.PasswordHash,
	)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return token, nil
}
