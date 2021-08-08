package controllers

import (
	"authserver/common"
	passwordhelpers "authserver/controllers/password_helpers"
	"authserver/models"
	"log"

	"github.com/google/uuid"
)

type CoreTokenController struct {
	PasswordHasher passwordhelpers.PasswordHasher
}

func (c CoreTokenController) CreateTokenFromPassword(CRUD TokenControllerCRUD, username string, password string, clientID uuid.UUID) (*models.AccessToken, common.OAuthCustomError) {
	//get the client
	client, rerr := parseClient(CRUD, clientID)
	if rerr.Type != common.ErrorTypeNone {
		return nil, rerr
	}

	//get the user
	user, err := CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, common.OAuthInternalError()
	}

	//check if user was found
	if user == nil {
		return nil, common.OAuthClientError("invalid_grant", "invalid username and/or password")
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, password)
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return nil, common.OAuthClientError("invalid_grant", "invalid username and/or password")
	}

	//create a new access token
	token := models.CreateNewAccessToken(user, client)

	//save the token
	err = CRUD.SaveAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error saving access token", err))
		return nil, common.OAuthInternalError()
	}

	return token, common.OAuthNoError()
}

func (c CoreTokenController) DeleteToken(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	//delete the token
	err := CRUD.DeleteAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error deleting access token", err))
		return common.InternalError()
	}

	//return success
	return common.NoError()
}

func (c CoreTokenController) DeleteAllOtherUserTokens(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	//delete the token
	err := CRUD.DeleteAllOtherUserTokens(token)
	if err != nil {
		log.Println(common.ChainError("error deleting all other user tokens", err))
		return common.InternalError()
	}

	//return success
	return common.NoError()
}
