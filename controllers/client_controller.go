package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CoreClientController struct{}

func (c CoreClientController) CreateClient(CRUD ClientControllerCRUD, name string) (*models.Client, common.CustomError) {
	//create the client model
	client := models.CreateNewClient(name)

	//validate the client
	verr := validateClient(client)
	if verr.Type != common.ErrorTypeNone {
		return nil, verr
	}

	//save the client
	err := CRUD.SaveClient(client)
	if err != nil {
		log.Println(common.ChainError("error saving client", err))
		return nil, common.InternalError()
	}

	return client, common.NoError()
}

func (c CoreClientController) UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError {
	//validate the client
	verr := validateClient(client)
	if verr.Type != common.ErrorTypeNone {
		return verr
	}

	//update the client
	res, err := CRUD.UpdateClient(client)
	if err != nil {
		log.Println(common.ChainError("error updating client", err))
		return common.InternalError()
	}

	//verify client was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("client with id %s not found", client.ID))
	}

	return common.NoError()
}

func (c CoreClientController) DeleteClient(CRUD ClientControllerCRUD, id uuid.UUID) common.CustomError {
	//delete the client with id
	res, err := CRUD.DeleteClient(id)
	if err != nil {
		log.Println(common.ChainError("error deleting client", err))
		return common.InternalError()
	}

	//verify client was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("client with id %s not found", id))
	}

	return common.NoError()
}

func validateClient(client *models.Client) common.CustomError {
	verr := client.Validate()
	if verr&models.ValidateClientEmptyName != 0 {
		return common.ClientError("client name cannot be empty")
	} else if verr&models.ValidateClientNameTooLong != 0 {
		return common.ClientError(fmt.Sprint("client name cannot be longer than ", models.ClientNameMaxLength, " characters"))
	}

	return common.NoError()
}
