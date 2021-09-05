package controllers

import (
	"authserver/common"
	"authserver/models"

	"github.com/google/uuid"
)

// Controllers encapsulates all other controller interfaces.
type Controllers interface {
	UserController
	ClientController
	UserRoleController
	AuthController
	SessionController
	TokenController
}

type CoreControllers struct {
	UserController
	ClientController
	UserRoleController
	AuthController
	SessionController
	TokenController
}

// UserControllerCRUD encapsulates the CRUD operations required by the UserController.
type UserControllerCRUD interface {
	models.UserCRUD
	models.SessionCRUD
}

type UserController interface {
	// CreateUser creates a new user with the given username and password.
	CreateUser(CRUD UserControllerCRUD, username string, password string) (*models.User, common.CustomError)

	// UpdateUserPassword updates the password for the given username.
	UpdateUserPassword(CRUD UserControllerCRUD, username string, oldPassword string, newPassword string) common.CustomError

	// DeleteUser deletes the user with given username.
	DeleteUser(CRUD UserControllerCRUD, username string) common.CustomError
}

// ClientControllerCRUD encapsulates the CRUD operations required by the ClientController.
type ClientControllerCRUD interface {
	models.ClientCRUD
	models.SessionCRUD
}

type ClientController interface {
	// CreateClient creates a new client with the given name and redirect url.
	CreateClient(CRUD ClientControllerCRUD, name string, redirectUrl string, tokenType int, keyUri string) (*models.Client, common.CustomError)

	// UpdateClient updates the given client.
	UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError

	// DeleteClient deletes the client with the given uid.
	DeleteClient(CRUD ClientControllerCRUD, uid uuid.UUID) common.CustomError
}

// UserRoleControllerCRUD encapsulates the CRUD operations required by the UserRoleController.
type UserRoleControllerCRUD interface {
	models.UserRoleCRUD
}

type UserRoleController interface {
	// UpdateUserRolesForClient updates the users' role for the provided client.
	UpdateUserRolesForClient(CRUD UserRoleControllerCRUD, clientUID uuid.UUID, roles []*models.UserRole) common.CustomError
}

// AuthControllerCRUD encapsulates the CRUD operations required by the AuthController.
type AuthControllerCRUD interface {
	models.UserCRUD
}

type AuthController interface {
	// AuthenticateUserWithPassword authenticates a user with their username and password.
	// Returns the user if authentication was successful, or nil if not.
	// Also returns any errors.
	AuthenticateUserWithPassword(CRUD AuthControllerCRUD, username string, password string) (*models.User, common.CustomError)
}

// SessionControllerCRUD encapsulates the CRUD operations required by the SessionController.
type SessionControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
	models.SessionCRUD
}

type SessionController interface {
	// CreateSession creates a new session by authorizing the user with a password.
	CreateSession(CRUD SessionControllerCRUD, username string, password string) (*models.Session, common.CustomError)

	// DeleteSession deletes the session with the given id.
	DeleteSession(CRUD SessionControllerCRUD, id uuid.UUID) common.CustomError

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username expect the one with the given id.
	DeleteAllOtherUserSessions(CRUD SessionControllerCRUD, username string, id uuid.UUID) common.CustomError
}

// TokenControllerCRUD encapsulates the CRUD operations required by the TokenController.
type TokenControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
}

type TokenController interface {
	// CreateTokenRedirectURL first authenticates using the username and password, then creates a signed JWT for the specified client.
	// The base-64 encoded token string is then appended to the client's redirect url.
	// Returns the url and any errors.
	CreateTokenRedirectURL(CRUD TokenControllerCRUD, clientId uuid.UUID, username string, password string) (string, common.CustomError)
}
