package controllers

import (
	"fmt"
	"log"

	"github.com/mhogar/amber/common"
	passwordhelpers "github.com/mhogar/amber/controllers/password_helpers"
	"github.com/mhogar/amber/models"
)

type CoreUserController struct {
	PasswordHasher            passwordhelpers.PasswordHasher
	PasswordCriteriaValidator passwordhelpers.PasswordCriteriaValidator
	AuthController            AuthController
	UserController            UserController
}

func (c CoreUserController) CreateUser(CRUD UserControllerCRUD, username string, password string, rank int) (*models.User, common.CustomError) {
	//create the user model
	user := models.CreateUser(username, rank, nil)

	//validate the user
	cerr := c.validateUser(user)
	if cerr.Type != common.ErrorTypeNone {
		return nil, cerr
	}

	//validate username is unique
	otherUser, err := CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, common.InternalError()
	}
	if otherUser != nil {
		return nil, common.ClientError("username is already in use")
	}

	//validate password meets criteria
	vperr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(password)
	if vperr.Status != passwordhelpers.ValidatePasswordCriteriaValid {
		log.Println(common.ChainError("error validating password criteria", vperr))
		return nil, common.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	user.PasswordHash, err = c.PasswordHasher.HashPassword(password)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return nil, common.InternalError()
	}

	//save the user
	err = CRUD.CreateUser(user)
	if err != nil {
		log.Println(common.ChainError("error saving user", err))
		return nil, common.InternalError()
	}

	return user, common.NoError()
}

func (CoreUserController) GetUsersWithLesserRank(CRUD UserControllerCRUD, rank int) ([]*models.User, common.CustomError) {
	//get the users
	users, err := CRUD.GetUsersWithLesserRank(rank)
	if err != nil {
		log.Println(common.ChainError("error getting users with lesser rank", err))
		return nil, common.InternalError()
	}

	return users, common.NoError()
}

func (c CoreUserController) UpdateUser(CRUD UserControllerCRUD, username string, rank int) (*models.User, common.CustomError) {
	//create the user model
	user := models.CreateUser(username, rank, nil)

	//validate the user
	cerr := c.validateUser(user)
	if cerr.Type != common.ErrorTypeNone {
		return nil, cerr
	}

	//update the user
	res, err := CRUD.UpdateUser(user)
	if err != nil {
		log.Println(common.ChainError("error updating user", err))
		return nil, common.InternalError()
	}

	//verify user was actually found
	if !res {
		return nil, common.ClientError(fmt.Sprintf("user with username %s not found", username))
	}

	return user, common.NoError()
}

func (c CoreUserController) UpdateUserPassword(CRUD UserControllerCRUD, username string, password string) common.CustomError {
	//validate password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(password)
	if verr.Status != passwordhelpers.ValidatePasswordCriteriaValid {
		log.Println(common.ChainError("error validating password criteria", verr))
		return common.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(password)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return common.InternalError()
	}

	//update the user's password (don't check result because we know the user already exists)
	_, err = CRUD.UpdateUserPassword(username, hash)
	if err != nil {
		log.Println(common.ChainError("error updating user password", err))
		return common.InternalError()
	}

	//return success
	return common.NoError()
}

func (c CoreUserController) UpdateUserPasswordWithAuth(CRUD UserControllerCRUD, username string, oldPassword string, newPassword string) common.CustomError {
	//authenticate user with their old password
	_, cerr := c.AuthController.AuthenticateUserWithPassword(CRUD, username, oldPassword)
	if cerr.Type == common.ErrorTypeClient {
		return common.ClientError("old password is incorrect")
	}
	if cerr.Type != common.ErrorTypeNone {
		return cerr
	}

	//update the user's password
	return c.UserController.UpdateUserPassword(CRUD, username, newPassword)
}

func (c CoreUserController) DeleteUser(CRUD UserControllerCRUD, username string) common.CustomError {
	//delete the user
	res, err := CRUD.DeleteUser(username)
	if err != nil {
		log.Println(common.ChainError("error deleting user", err))
		return common.InternalError()
	}

	//verify user was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("user with username %s not found", username))
	}

	//return success
	return common.NoError()
}

func (CoreUserController) VerifyUserRank(CRUD UserControllerCRUD, username string, rank int) (bool, common.CustomError) {
	//get the requested user
	user, err := CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return false, common.InternalError()
	}

	//verify user exists
	if user == nil {
		return false, common.ClientError("the requested user was not found")
	}

	//verify the rank
	return user.Rank < rank, common.NoError()
}

func (CoreUserController) validateUser(user *models.User) common.CustomError {
	verr := user.Validate()

	if verr&models.ValidateUserEmptyUsername != 0 {
		return common.ClientError("username cannot be empty")
	}
	if verr&models.ValidateUserUsernameTooLong != 0 {
		return common.ClientError(fmt.Sprint("username cannot be longer than ", models.UserUsernameMaxLength, " characters"))
	}
	if verr&models.ValidateUserInvalidRank != 0 {
		return common.ClientError("rank is invalid")
	}

	return common.NoError()
}
