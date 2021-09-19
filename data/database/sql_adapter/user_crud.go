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

func (crud *SQLCRUD) CreateUser(user *models.User) error {
	//validate the user model
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return errors.New(fmt.Sprint("error validating user model:", verr))
	}

	if user.PasswordHash == nil {
		return errors.New("password hash cannot be nil")
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserScript(),
		user.Username, user.Rank, user.PasswordHash,
	)
	cancel()

	if err != nil {
		return common.ChainError("error executing create user statement", err)
	}

	return nil
}

func (crud *SQLCRUD) GetUsersWithLesserRank(rank int) ([]*models.User, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUsersWithLesserRankScript(), rank)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get users with lesser rank query", err)
	}
	defer rows.Close()

	//read the data
	users := []*models.User{}
	for {
		user, err := readUserData(rows)
		if err != nil {
			return nil, err
		}

		if user == nil {
			break
		}
		users = append(users, user)
	}
	return users, nil
}

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

func (crud *SQLCRUD) UpdateUser(user *models.User) (bool, error) {
	//validate the user model
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return false, errors.New(fmt.Sprint("error validating user model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateUserScript(),
		user.Username, user.Rank,
	)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update user statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

func (crud *SQLCRUD) UpdateUserPassword(username string, hash []byte) (bool, error) {
	if hash == nil {
		return false, errors.New("password hash cannot be nil")
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateUserPasswordScript(),
		username, hash,
	)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update user password statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

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
	err := rows.Scan(
		&user.Username, &user.Rank, &user.PasswordHash,
	)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return user, nil
}
