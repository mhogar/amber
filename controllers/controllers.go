package controllers

import (
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"

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
	// CreateUser creates a new user with the given username, password, and rank.
	// Returns the user model and any errors.
	CreateUser(CRUD UserControllerCRUD, username string, password string, rank int) (*models.User, common.CustomError)

	// GetUsersWithLesserRank gets all users with a rank less than the provided one.
	// Returns the user models and any errors.
	GetUsersWithLesserRank(CRUD UserControllerCRUD, rank int) ([]*models.User, common.CustomError)

	// UpdateUser updates the fields of the user for the given username.
	// Returns the user model and any errors.
	UpdateUser(CRUD UserControllerCRUD, username string, rank int) (*models.User, common.CustomError)

	// UpdateUserPassword updates the password for the user with the given username.
	// Returns any errors.
	UpdateUserPassword(CRUD UserControllerCRUD, username string, password string) common.CustomError

	// UpdateUserPasswordWithAuth authenticates the user and updates their password.
	// Returns any errors.
	UpdateUserPasswordWithAuth(CRUD UserControllerCRUD, username string, oldPassword string, newPassword string) common.CustomError

	// DeleteUser deletes the user with given username.
	// Returns any errors.
	DeleteUser(CRUD UserControllerCRUD, username string) common.CustomError

	// VerifyUserRank verifies the user with given username has a rank less than the provided rank.
	// Returns result and any errors.
	VerifyUserRank(CRUD UserControllerCRUD, username string, rank int) (bool, common.CustomError)
}

// ClientControllerCRUD encapsulates the CRUD operations required by the ClientController.
type ClientControllerCRUD interface {
	models.ClientCRUD
	models.SessionCRUD
}

type ClientController interface {
	// CreateClient creates a new client using the provided model.
	// Returns any errors.
	CreateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError

	// GetClients gets all clients.
	// Returns the client models any errors.
	GetClients(CRUD ClientControllerCRUD) ([]*models.Client, common.CustomError)

	// UpdateClient updates the given client.
	// Returns any errors.
	UpdateClient(CRUD ClientControllerCRUD, client *models.Client) common.CustomError

	// DeleteClient deletes the client with the given uid.
	// Returns any errors.
	DeleteClient(CRUD ClientControllerCRUD, uid uuid.UUID) common.CustomError
}

// UserRoleControllerCRUD encapsulates the CRUD operations required by the UserRoleController.
type UserRoleControllerCRUD interface {
	models.UserRoleCRUD
}

type UserRoleController interface {
	// CreateUserRole creates a new user-role using the provided model.
	// Returns any errors.
	CreateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError

	// GetUserRolesByClientUID gets the user-roles with the provided client uid and with a rank less than the provided rank.
	// Returns the user-role models and any errors.
	GetUserRolesWithLesserRankByClientUID(CRUD UserRoleControllerCRUD, clientUID uuid.UUID, rank int) ([]*models.UserRole, common.CustomError)

	// UpdateUserRole updates the given user-role.
	// Returns any errors.
	UpdateUserRole(CRUD UserRoleControllerCRUD, role *models.UserRole) common.CustomError

	// DeleteUserRole deletes the user-role with the given client uid and username.
	// Returns any errors.
	DeleteUserRole(CRUD UserRoleControllerCRUD, clientUID uuid.UUID, username string) common.CustomError
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
	models.SessionCRUD
}

type SessionController interface {
	// CreateSession creates a new session by authorizing the user with a password.
	// Returns the session model and any errors.
	CreateSession(CRUD SessionControllerCRUD, username string, password string) (*models.Session, common.CustomError)

	// DeleteSession deletes the session with the given id.
	// Returns any errors.
	DeleteSession(CRUD SessionControllerCRUD, id uuid.UUID) common.CustomError

	// DeleteAllUserSessions deletes all of the sessions for the given username.
	// Returns any errors.
	DeleteAllUserSessions(CRUD SessionControllerCRUD, username string) common.CustomError

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username expect the one with the given id.
	// Returns any errors.
	DeleteAllOtherUserSessions(CRUD SessionControllerCRUD, username string, id uuid.UUID) common.CustomError
}

// TokenControllerCRUD encapsulates the CRUD operations required by the TokenController.
type TokenControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
	models.UserRoleCRUD
}

type TokenController interface {
	// CreateTokenRedirectURL first authenticates using the username and password, then creates a signed JWT for the specified client.
	// The base-64 encoded token string is then appended to the client's redirect url.
	// Returns the url and any errors.
	CreateTokenRedirectURL(CRUD TokenControllerCRUD, clientId uuid.UUID, username string, password string) (string, common.CustomError)
}
