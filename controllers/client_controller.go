package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"
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
	err := CRUD.CreateClient(client)
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
		return common.ClientError(fmt.Sprintf("client with id %s not found", client.UID))
	}

	return common.NoError()
}

func (c CoreClientController) DeleteClient(CRUD ClientControllerCRUD, id int16) common.CustomError {
	//delete the client
	res, err := CRUD.DeleteClient(id)
	if err != nil {
		log.Println(common.ChainError("error deleting client", err))
		return common.InternalError()
	}

	//verify client was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("client with id %d not found", id))
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
