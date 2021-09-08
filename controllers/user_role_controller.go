package controllers

import (
	"authserver/common"
	"authserver/models"

	"github.com/google/uuid"
)

type CoreUserRoleController struct{}

func (c CoreUserRoleController) CreateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	//TODO: verify the user does not already have a role for the client

	return common.NoError()
}

func (c CoreUserRoleController) UpdateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	return common.NoError()
}

func (c CoreUserRoleController) DeleteUserRole(CRUD UserRoleControllerCRUD, username string, clientUID uuid.UUID) common.CustomError {
	return common.NoError()
}
