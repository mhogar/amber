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
func (crud *SQLCRUD) SaveSession(session *models.Session) error {
	verr := session.Validate()
	if verr != models.ValidateSessionValid {
		return errors.New(fmt.Sprint("error validating session model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveSessionScript(), session.Token, session.Username)
	cancel()

	if err != nil {
		return common.ChainError("error executing save session statement", err)
	}

	return nil
}

// GetSessionByToken gets the row in the session table with the matching token, and creates a new session model using its data.
// Returns the model and any errors.
func (crud *SQLCRUD) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetSessionByTokenScript(), token)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get session by token query", err)
	}
	defer rows.Close()

	return readSessionData(rows)
}

// DeleteSession deletes the row in the session table with the matching token.
// Returns result of whether the session was found, and any errors.
func (crud *SQLCRUD) DeleteSession(token uuid.UUID) (bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteSessionScript(), token)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing delete session statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

// DeleteAllOtherUserSessions deletes all the rows in the session table with the matching user token, and not the session token.
// Returns any errors.
func (crud *SQLCRUD) DeleteAllOtherUserSessions(username string, token uuid.UUID) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteAllOtherUserSessionsScript(), token, username)
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

	session := &models.Session{}

	//get the result
	err := rows.Scan(
		&session.Token, &session.Username, &session.Rank,
	)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return session, nil
}
