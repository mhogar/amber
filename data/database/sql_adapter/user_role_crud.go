package sqladapter

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
)

// CreateUserRoleTable creates the user-role table in the database.
// Returns any errors.
func (crud *SQLCRUD) CreateUserRoleTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserRoleTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create user-role table script", err)
	}

	return err
}

// DropUserRoleTable drops the user-role table from the database.
// Returns any errors.
func (crud *SQLCRUD) DropUserRoleTable() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DropUserRoleTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing drop user-role table script", err)
	}

	return err
}

func (crud *SQLCRUD) CreateUserRole(role *models.UserRole) error {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateUserRoleScript(),
		role.ClientUID, role.Username, role.Role,
	)
	defer cancel()

	if err != nil {
		return common.ChainError("error executing create user role statement", err)
	}

	return nil
}

func (crud *SQLCRUD) GetUserRolesWithLesserRankByClientUID(uid uuid.UUID, rank int) ([]*models.UserRole, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserRolesWithLesserRankByClientUIDScript(), uid, rank)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user roles with lesser rank by client uid query", err)
	}
	defer rows.Close()

	//read the data
	roles := []*models.UserRole{}
	for {
		role, err := readUserRoleData(rows)
		if err != nil {
			return nil, err
		}

		if role == nil {
			break
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (crud *SQLCRUD) GetUserRoleByClientUIDAndUsername(clientUID uuid.UUID, username string) (*models.UserRole, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserRoleByClientUIDAndUsernameScript(),
		clientUID, username,
	)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user-roles by client uid and username query", err)
	}
	defer rows.Close()

	return readUserRoleData(rows)
}

func (crud *SQLCRUD) UpdateUserRole(role *models.UserRole) (bool, error) {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return false, errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.UpdateUserRoleScript(),
		role.ClientUID, role.Username, role.Role,
	)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing update user role statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

func (crud *SQLCRUD) DeleteUserRole(clientUID uuid.UUID, username string) (bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	res, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteUserRoleScript(),
		clientUID, username,
	)
	cancel()

	if err != nil {
		return false, common.ChainError("error executing delete user role statement", err)
	}

	count, _ := res.RowsAffected()
	return count > 0, nil
}

func readUserRoleData(rows *sql.Rows) (*models.UserRole, error) {
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
	userRole := &models.UserRole{}
	err := rows.Scan(
		&userRole.ClientUID, &userRole.Username, &userRole.Role,
	)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return userRole, nil
}
