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
	AuthController
	SessionController
}

type CoreControllers struct {
	UserController
	ClientController
	AuthController
	SessionController
}

// UserControllerCRUD encapsulates the CRUD operations required by the UserController.
type UserControllerCRUD interface {
	models.UserCRUD
	models.SessionCRUD
}

type UserController interface {
	// CreateUser creates a new user with the given username and password.
	CreateUser(CRUD UserControllerCRUD, username string, password string) (*models.User, common.CustomError)

	// UpdateUserPassword updates the given user's password.
	UpdateUserPassword(CRUD UserControllerCRUD, user *models.User, oldPassword string, newPassword string) common.CustomError

	// DeleteUser deletes the user with given username.
	DeleteUser(CRUD UserControllerCRUD, username string) common.CustomError
}

// ClientControllerCRUD encapsulates the CRUD operations required by the ClientController.
type ClientControllerCRUD interface {
	models.ClientCRUD
	models.SessionCRUD
}

type ClientController interface {
	// CreateClient creates a new client with the given name.
	CreateClient(CRUD ClientControllerCRUD, name string) (*models.Client, common.CustomError)

	// UpdateClient updates the given client.
	UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError

	// DeleteClient deletes the client with the given uid.
	DeleteClient(CRUD ClientControllerCRUD, uid uuid.UUID) common.CustomError
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

	// DeleteSession deletes the session.
	DeleteSession(CRUD SessionControllerCRUD, session *models.Session) common.CustomError

	// DeleteSession deletes all of the user's sessions accept for the provided one.
	DeleteAllOtherUserSessions(CRUD SessionControllerCRUD, session *models.Session) common.CustomError
}
