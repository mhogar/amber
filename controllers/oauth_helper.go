package controllers

import (
	"authserver/common"
	"authserver/models"
	"log"

	"github.com/google/uuid"
)

func parseClient(clientCRUD models.ClientCRUD, clientID uuid.UUID) (*models.Client, common.OAuthCustomError) {
	//get the client
	client, err := clientCRUD.GetClientByID(clientID)
	if err != nil {
		log.Println(common.ChainError("error getting client by id", err))
		return nil, common.OAuthInternalError()
	}

	//check client was found
	if client == nil {
		return nil, common.OAuthClientError("invalid_client", "client with id not found")
	}

	return client, common.OAuthNoError()
}

func parseScope(scopeCRUD models.ScopeCRUD, name string) (*models.Scope, common.OAuthCustomError) {
	//get the scope
	scope, err := scopeCRUD.GetScopeByName(name)
	if err != nil {
		log.Println(common.ChainError("error getting scope by name", err))
		return nil, common.OAuthInternalError()
	}

	if scope == nil {
		return nil, common.OAuthClientError("invalid_scope", "scope with name not found")
	}

	return scope, common.OAuthNoError()
}
