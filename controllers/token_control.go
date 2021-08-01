package controllers

import (
	"authserver/common"
	passwordhelpers "authserver/controllers/password_helpers"
	"authserver/models"
	"log"

	"github.com/google/uuid"
)

// TokenControl handles requests to "/token" endpoints
type TokenControl struct {
	PasswordHasher passwordhelpers.PasswordHasher
}

// PostToken handles POST requests to "/token"
func (c TokenControl) CreateTokenFromPassword(CRUD TokenControllerCRUD, username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, common.OAuthCustomError) {
	//get the client
	client, rerr := parseClient(CRUD, clientID)
	if rerr.Type != common.ErrorTypeNone {
		return nil, rerr
	}

	//get the scope
	scope, rerr := parseScope(CRUD, scopeName)
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
	token := models.CreateNewAccessToken(user, client, scope)

	//save the token
	err = CRUD.SaveAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error saving access token", err))
		return nil, common.OAuthInternalError()
	}

	return token, common.OAuthNoError()
}

// DeleteToken deletes the access token.
func (c TokenControl) DeleteToken(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	//delete the token
	err := CRUD.DeleteAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error deleting access token", err))
		return common.InternalError()
	}

	//return success
	return common.NoError()
}

// DeleteToken deletes all of the user's tokens accept for the provided one.
func (c TokenControl) DeleteAllOtherUserTokens(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	//delete the token
	err := CRUD.DeleteAllOtherUserTokens(token)
	if err != nil {
		log.Println(common.ChainError("error deleting all other user tokens", err))
		return common.InternalError()
	}

	//return success
	return common.NoError()
}
