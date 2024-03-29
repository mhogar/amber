package controllers

import (
	"fmt"
	"log"
	"net/url"

	"github.com/mhogar/amber/common"
	jwthelpers "github.com/mhogar/amber/controllers/jwt_helpers"

	"github.com/google/uuid"
)

type CoreTokenController struct {
	AuthController       AuthController
	TokenFactorySelector jwthelpers.TokenFactorySelector
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
	if cerr.Type == common.ErrorTypeClient {
		return "", common.ClientError("invalid username and/or password, or user is not assigned to the client")
	}
	if cerr.Type != common.ErrorTypeNone {
		return "", cerr
	}

	//get the user's role
	role, err := CRUD.GetUserRoleByClientUIDAndUsername(clientUID, username)
	if err != nil {
		log.Println(common.ChainError("error getting user role", err))
		return "", common.InternalError()
	}

	//verify role exists
	if role == nil {
		return "", common.ClientError("invalid username and/or password, or user is not assigned to the client")
	}

	//choose the token factory (in practice a factory should always be found since the client model validates the token type when saving)
	tf := c.TokenFactorySelector.Select(client.TokenType)
	if tf == nil {
		log.Println(fmt.Sprintf("token factory for token type %d not found", client.TokenType))
		return "", common.InternalError()
	}

	//create the token
	token, err := tf.CreateToken(client.KeyUri, clientUID, username, role.Role)
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
