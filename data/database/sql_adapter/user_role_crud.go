package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"database/sql"
	"errors"
	"fmt"

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

func (crud *SQLCRUD) GetUserRolesForClient(clientUID uuid.UUID) ([]*models.UserRole, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserRolesForClientScript(), clientUID)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user roles for client query", err)
	}
	defer rows.Close()

	//read the data
	userRoles := make([]*models.UserRole, 0)
	for {
		userRole, err := readUserRoleData(rows)
		if err != nil {
			return nil, err
		}

		if userRole == nil {
			break
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (crud *SQLCRUD) GetUserRoleForClient(clientUID uuid.UUID, username string) (*models.UserRole, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetUserRolesForClientScript(), clientUID, username)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get user role for client query", err)
	}
	defer rows.Close()

	return readUserRoleData(rows)
}

func (crud *SQLCRUD) UpdateUserRoles(clientUID uuid.UUID, roles []*models.UserRole) error {
	//validate the models
	for _, role := range roles {
		verr := role.Validate()
		if verr != models.ValidateUserRoleValid {
			return errors.New(fmt.Sprint("error validating user-role model:", verr))
		}
	}

	//-- delete all existing roles first --
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteUserRolesForClientScript(), clientUID)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete user roles for client statement", err)
	}

	//-- add new roles --
	for _, role := range roles {
		ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
		_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.AddUserRoleForClientScript(), clientUID, role.Username, role.Role)
		cancel()

		if err != nil {
			return common.ChainError("error executing add user role for client statement", err)
		}
	}

	return nil
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
	err := rows.Scan(&userRole.Username, &userRole.Role)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return userRole, nil
}
