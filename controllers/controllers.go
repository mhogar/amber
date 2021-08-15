package controllers

import (
	"authserver/common"
	"authserver/models"

	"github.com/google/uuid"
)

// Controllers encapsulates all other controller interfaces
type Controllers interface {
	UserController
	ClientController
	TokenController
}

type CoreControllers struct {
	CoreUserController
	CoreClientController
	CoreTokenController
}

// UserControllerCRUD encapsulates the CRUD operations required by the UserController
type UserControllerCRUD interface {
	models.UserCRUD
	models.AccessTokenCRUD
}

type UserController interface {
	// CreateUser creates a new user with the given username and password
	CreateUser(CRUD UserControllerCRUD, username string, password string) (*models.User, common.CustomError)

	// UpdateUserPassword updates the given user's password
	UpdateUserPassword(CRUD UserControllerCRUD, user *models.User, oldPassword string, newPassword string) common.CustomError

	// DeleteUser deletes the given id
	DeleteUser(CRUD UserControllerCRUD, id int32) common.CustomError
}

// ClientControllerCRUD encapsulates the CRUD operations required by the ClientController
type ClientControllerCRUD interface {
	models.ClientCRUD
	models.AccessTokenCRUD
}

type ClientController interface {
	// CreateClient creates a new client with the given name
	CreateClient(CRUD ClientControllerCRUD, name string) (*models.Client, common.CustomError)

	// UpdateClient updates the given client
	UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError

	// DeleteClient deletes the given id
	DeleteClient(CRUD ClientControllerCRUD, id int16) common.CustomError
}

// TokenControllerCRUD encapsulates the CRUD operations required by the TokenController
type TokenControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
	models.AccessTokenCRUD
}

type TokenController interface {
	// CreateTokenFromPassword creates a new access token, authenticating using a password
	CreateTokenFromPassword(CRUD TokenControllerCRUD, username string, password string, clientID uuid.UUID) (*models.AccessToken, common.OAuthCustomError)

	// DeleteToken deletes the access token
	DeleteToken(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError

	// DeleteToken deletes all of the user's tokens accept for the provided one
	DeleteAllOtherUserTokens(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError
}
