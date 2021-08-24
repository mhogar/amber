package controllers

import (
	"authserver/common"
	jwthelpers "authserver/controllers/jwt_helpers"
	"fmt"
	"log"
	"net/url"

	"github.com/google/uuid"
)

type CoreTokenController struct {
	AuthController AuthController
	TokenFactory   jwthelpers.TokenFactory
}

func (c CoreTokenController) CreateTokenRedirectURL(CRUD TokenControllerCRUD, clientUID uuid.UUID, username string, password string) (string, common.CustomError) {
	//get the requested client
	client, err := CRUD.GetClientByUID(clientUID)
	if err != nil {
		log.Println(common.ChainError("error getting client by uid", err))
		return "", common.InternalError()
	}

	//verify client exists
	if client == nil {
		return "", common.ClientError(fmt.Sprintf("client with id %s not found", clientUID.String()))
	}

	//authenticate the user
	_, cerr := c.AuthController.AuthenticateUserWithPassword(CRUD, username, password)
	if cerr.Type != common.ErrorTypeNone {
		return "", cerr
	}

	//create the token
	token, err := c.TokenFactory.CreateToken(username)
	if err != nil {
		log.Println(common.ChainError("error creating token", err))
		return "", common.InternalError()
	}

	//parse the redirect url (in practice this should always succeed since the client model validates the url when saving)
	url, err := url.Parse(client.RedirectUrl)
	if err != nil {
		log.Println(common.ChainError("error parsing redirect url", err))
		return "", common.InternalError()
	}

	//set the token as a query parameter
	q := url.Query()
	q.Set("token", token)
	url.RawQuery = q.Encode()

	return url.String(), common.NoError()
}
