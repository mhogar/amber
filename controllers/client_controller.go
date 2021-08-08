package controllers

import (
	"authserver/common"
	"authserver/models"

	"github.com/google/uuid"
)

type CoreClientController struct{}

func (c CoreClientController) CreateClient(CRUD ClientControllerCRUD, name string) (*models.Client, common.CustomError) {
	return nil, common.NoError()
}

func (c CoreClientController) UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError {
	return common.NoError()
}

func (c CoreClientController) DeleteClient(CRUD ClientControllerCRUD, id uuid.UUID) common.CustomError {
	return common.NoError()
}
