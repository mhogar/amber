package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CoreUserRoleController struct{}

func (c CoreUserRoleController) CreateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	//validate the model
	cerr := c.validateUserRole(role)
	if cerr.Type != common.ErrorTypeNone {
		return cerr
	}

	//verify the user does not already have a role for the client
	existingRole, err := CRUD.GetUserRoleByClientUIDAndUsername(role.ClientUID, role.Username)
	if err != nil {
		log.Println("error getting user-role by username and client uid", err)
		return common.InternalError()
	}
	if existingRole != nil {
		return common.ClientError("the user already has a role for the client")
	}

	//create the user-role
	err = CRUD.CreateUserRole(role)
	if err != nil {
		log.Println("error creating user-role", err)
		return common.InternalError()
	}

	return common.NoError()
}

func (c CoreUserRoleController) GetUserRolesByClientUID(CRUD UserRoleControllerCRUD, clientUID uuid.UUID) ([]*models.UserRole, common.CustomError) {
	//get the roles
	roles, err := CRUD.GetUserRolesByClientUID(clientUID)
	if err != nil {
		log.Println("error getting user-roles by client uid", err)
		return nil, common.InternalError()
	}

	return roles, common.NoError()
}

func (c CoreUserRoleController) UpdateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	//validate the model
	cerr := c.validateUserRole(role)
	if cerr.Type != common.ErrorTypeNone {
		return cerr
	}

	//update the user-role
	res, err := CRUD.UpdateUserRole(role)
	if err != nil {
		log.Println("error creating user-role", err)
		return common.InternalError()
	}

	//verify user-role was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("no role found for user %s and client %s", role.Username, role.ClientUID.String()))
	}

	return common.NoError()
}

func (c CoreUserRoleController) DeleteUserRole(CRUD UserRoleControllerCRUD, username string, clientUID uuid.UUID) common.CustomError {
	//delete the user-role
	res, err := CRUD.DeleteUserRole(username, clientUID)
	if err != nil {
		log.Println("error deleting user-role", err)
		return common.InternalError()
	}

	//verify user-role was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("no role found for user %s and client %s", username, clientUID.String()))
	}

	return common.NoError()
}

func (CoreUserRoleController) validateUserRole(role *models.UserRole) common.CustomError {
	verr := role.Validate()

	if verr&models.ValidateUserRoleEmptyRole != 0 {
		return common.ClientError("role cannot be empty")
	} else if verr&models.ValidateUserRoleRoleTooLong != 0 {
		return common.ClientError(fmt.Sprint("role cannot be longer than ", models.UserRoleRoleMaxLength, " characters"))
	}

	return common.NoError()
}
