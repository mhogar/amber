package controllers

import (
	"authserver/common"
	"authserver/models"

	"github.com/google/uuid"
)

// Controllers encapsulates all other controller interfaces
type Controllers interface {
	UserController
	TokenController
}

type CoreControllers struct {
	CoreUserController
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

	// DeleteUser deletes the given user
	DeleteUser(CRUD UserControllerCRUD, user *models.User) common.CustomError

	// UpdateUserPassword updates the given user's password
	UpdateUserPassword(CRUD UserControllerCRUD, user *models.User, oldPassword string, newPassword string) common.CustomError
}

// TokenControllerCRUD encapsulates the CRUD operations required by the TokenController
type TokenControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
	models.ScopeCRUD
	models.AccessTokenCRUD
}

type TokenController interface {
	// CreateTokenFromPassword creates a new access token, authenticating using a password
	CreateTokenFromPassword(CRUD TokenControllerCRUD, username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, common.OAuthCustomError)

	// DeleteToken deletes the access token
	DeleteToken(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError

	// DeleteToken deletes all of the user's tokens accept for the provided one
	DeleteAllOtherUserTokens(CRUD TokenControllerCRUD, token *models.AccessToken) common.CustomError
}
