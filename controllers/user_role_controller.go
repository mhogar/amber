package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CoreUserRoleController struct{}

func (c CoreUserRoleController) UpdateUserRolesForClient(CRUD UserRoleControllerCRUD, clientUID uuid.UUID, roles []*models.UserRole) common.CustomError {
	//validate the roles
	for _, role := range roles {
		cerr := c.validateUserRole(role)
		if cerr.Type != common.ErrorTypeNone {
			return cerr
		}
	}

	//update the users' roles
	err := CRUD.UpdateUserRolesForClient(clientUID, roles)
	if err != nil {
		log.Println(common.ChainError("error updating user roles for client", err))
		return common.InternalError()
	}

	return common.NoError()
}

func (CoreUserRoleController) validateUserRole(role *models.UserRole) common.CustomError {
	verr := role.Validate()

	if verr&models.ValidateUserRoleEmptyUsername != 0 {
		return common.ClientError("username cannot be empty")
	} else if verr&models.ValidateUserRoleUsernameTooLong != 0 {
		return common.ClientError(fmt.Sprint("username cannot be longer than ", models.UserUsernameMaxLength, " characters"))
	} else if verr&models.ValidateUserRoleEmptyRole != 0 {
		return common.ClientError("role cannot be empty")
	} else if verr&models.ValidateUserRoleRoleTooLong != 0 {
		return common.ClientError(fmt.Sprint("role cannot be longer than ", models.UserRoleRoleMaxLength, " characters"))
	}

	return common.NoError()
}
