package controllers

import (
	"fmt"
	"log"

	"authserver/common"
	passwordhelpers "authserver/controllers/password_helpers"
	"authserver/models"
)

type CoreUserController struct {
	PasswordHasher            passwordhelpers.PasswordHasher
	PasswordCriteriaValidator passwordhelpers.PasswordCriteriaValidator
}

func (c CoreUserController) CreateUser(CRUD UserControllerCRUD, username string, password string) (*models.User, common.CustomError) {
	//create the user model
	user := models.CreateNewUser(username, nil)

	//validate the username
	verr := user.Validate()
	if verr&models.ValidateUserEmptyUsername != 0 {
		return nil, common.ClientError("username cannot be empty")
	} else if verr&models.ValidateUserUsernameTooLong != 0 {
		return nil, common.ClientError(fmt.Sprint("username cannot be longer than ", models.UserUsernameMaxLength, " characters"))
	}

	//validate username is unique
	otherUser, err := CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, common.InternalError()
	}
	if otherUser != nil {
		return nil, common.ClientError("error creating user")
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

func (c CoreUserController) DeleteUser(CRUD UserControllerCRUD, id int32) common.CustomError {
	//delete the user
	res, err := CRUD.DeleteUser(id)
	if err != nil {
		log.Println(common.ChainError("error deleting user", err))
		return common.InternalError()
	}

	//verify user was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("user with id %d not found", id))
	}

	//return success
	return common.NoError()
}

func (c CoreUserController) UpdateUserPassword(CRUD UserControllerCRUD, user *models.User, oldPassword string, newPassword string) common.CustomError {
	//validate old password
	err := c.PasswordHasher.ComparePasswords(user.PasswordHash, oldPassword)
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return common.ClientError("old password is invalid")
	}

	//validate new password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(newPassword)
	if verr.Status != passwordhelpers.ValidatePasswordCriteriaValid {
		log.Println(common.ChainError("error validating password criteria", verr))
		return common.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(newPassword)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return common.InternalError()
	}

	//update the user
	user.PasswordHash = hash
	res, err := CRUD.UpdateUser(user)
	if err != nil {
		log.Println(common.ChainError("error updating user", err))
		return common.InternalError()
	}

	//verify user was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("user %s not found", user.Username))
	}

	//return success
	return common.NoError()
}
