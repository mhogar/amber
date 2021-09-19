package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CoreClientController struct{}

func (c CoreClientController) CreateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError {
	//validate the client
	verr := c.validateClient(client)
	if verr.Type != common.ErrorTypeNone {
		return verr
	}

	//save the client
	err := CRUD.CreateClient(client)
	if err != nil {
		log.Println(common.ChainError("error saving client", err))
		return common.InternalError()
	}

	return common.NoError()
}

func (CoreClientController) GetClients(CRUD ClientControllerCRUD) ([]*models.Client, common.CustomError) {
	//get the clients
	clients, err := CRUD.GetClients()
	if err != nil {
		log.Println(common.ChainError("error getting clients", err))
		return nil, common.InternalError()
	}

	return clients, common.NoError()
}

func (c CoreClientController) UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError {
	//validate the client
	verr := c.validateClient(client)
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

func (CoreClientController) DeleteClient(CRUD ClientControllerCRUD, uid uuid.UUID) common.CustomError {
	//delete the client
	res, err := CRUD.DeleteClient(uid)
	if err != nil {
		log.Println(common.ChainError("error deleting client", err))
		return common.InternalError()
	}

	//verify client was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("client with id %s not found", uid.String()))
	}

	return common.NoError()
}

func (CoreClientController) validateClient(client *models.Client) common.CustomError {
	verr := client.Validate()

	if verr&models.ValidateClientEmptyName != 0 {
		return common.ClientError("client name cannot be empty")
	}
	if verr&models.ValidateClientNameTooLong != 0 {
		return common.ClientError(fmt.Sprint("client name cannot be longer than ", models.ClientNameMaxLength, " characters"))
	}
	if verr&models.ValidateClientEmptyRedirectUrl != 0 {
		return common.ClientError("client redirect url cannot be empty")
	}
	if verr&models.ValidateClientRedirectUrlTooLong != 0 {
		return common.ClientError(fmt.Sprint("client redirect url cannot be longer than ", models.ClientRedirectUrlMaxLength, " characters"))
	}
	if verr&models.ValidateClientInvalidRedirectUrl != 0 {
		return common.ClientError("client redirect url is an invalid url")
	}
	if verr&models.ValidateClientInvalidTokenType != 0 {
		return common.ClientError("client token type is invalid")
	}
	if verr&models.ValidateClientEmptyKeyUri != 0 {
		return common.ClientError("client key uri cannot be empty")
	}
	if verr&models.ValidateClientKeyUriTooLong != 0 {
		return common.ClientError(fmt.Sprint("client key uri cannot be longer than ", models.ClientKeyUriMaxLength, " characters"))
	}

	return common.NoError()
}
