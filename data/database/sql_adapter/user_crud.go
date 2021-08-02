package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CreateUserTable creates the user table in the database
// Returns any errors
func (crud *SQLCRUD) CreateUserTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create user table script", err)
	}

	return err
}

// DropUserTable drops the user table from the database
// Returns any errors
func (crud *SQLCRUD) DropUserTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropUserTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop user table script", err)
	}

	return err
}

// SaveUser validates the user model is valid and inserts a new row into the user table
// Returns any errors
func (crud *SQLCRUD) SaveUser(user *models.User) error {
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return errors.New(fmt.Sprint("error validating user model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveUserScript(),
		user.ID, user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return common.ChainError("error executing save user statement", err)
	}

	return nil
}

// GetUserByID gets the row in the user table with the matching id, and creates a new user model using its data
// Returns the model and any errors
func (crud *SQLCRUD) GetUserByID(ID uuid.UUID) (*models.User, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserByIdScript(), ID)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user by id query", err)
	}
	defer rows.Close()

	return readUserData(rows)
}

// GetUserByUsername gets the row in the user table with the matching username, and creates a new user model using its data
// Returns the model and any errors
func (crud *SQLCRUD) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserByUsernameScript(), username)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user by username query", err)
	}
	defer rows.Close()

	return readUserData(rows)
}

// UpdateUser validates the user model is valid and updates the row in the user table with the matching id
// Returns any errors
func (crud *SQLCRUD) UpdateUser(user *models.User) error {
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return errors.New(fmt.Sprint("error validating user model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateUserScript(),
		user.ID, user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return common.ChainError("error executing update user statement", err)
	}

	return nil
}

// DeleteUser deletes the row in the user table with the matching id
// Returns any errors
func (crud *SQLCRUD) DeleteUser(user *models.User) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteUserScript(), user.ID)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete user statement", err)
	}

	return nil
}

func readUserData(rows *sql.Rows) (*models.User, error) {
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
	user := &models.User{}
	err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return user, nil
}
