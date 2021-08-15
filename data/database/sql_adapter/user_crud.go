package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"
)

// CreateUserTable creates the user table in the database.
// Returns any errors.
func (crud *SQLCRUD) CreateUserTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create user table script", err)
	}

	return err
}

// DropUserTable drops the user table from the database.
// Returns any errors.
func (crud *SQLCRUD) DropUserTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropUserTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop user table script", err)
	}

	return err
}

// CreateUser validates the user model is valid and inserts a new row into the user table.
// Updates the model with the new inserted id and returns any errors.
func (crud *SQLCRUD) CreateUser(user *models.User) error {
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return errors.New(fmt.Sprint("error validating user model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserScript(),
		user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return common.ChainError("error executing create user statement", err)
	}

	return nil
}

// GetUserByUsername gets the row in the user table with the matching username, and creates a new user model using its data.
// Returns the model and any errors.
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

// UpdateUser validates the user model is valid and updates the row in the user table with the matching id.
// Returns result of whether the user was found, and any errors.
func (crud *SQLCRUD) UpdateUser(user *models.User) (bool, error) {
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return false, errors.New(fmt.Sprint("error validating user model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateUserScript(),
		user.Username, user.PasswordHash)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update user statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

// DeleteUser deletes the row in the user table with the matching id.
// Returns result of whether the user was found, and any errors.
func (crud *SQLCRUD) DeleteUser(username string) (bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteUserScript(), username)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing delete user statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
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
	err := rows.Scan(&user.Username, &user.PasswordHash)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return user, nil
}
