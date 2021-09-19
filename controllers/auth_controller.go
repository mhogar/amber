package controllers

import (
	"log"

	"github.com/mhogar/amber/common"
	passwordhelpers "github.com/mhogar/amber/controllers/password_helpers"
	"github.com/mhogar/amber/models"
)

type CoreAuthController struct {
	PasswordHasher passwordhelpers.PasswordHasher
}

func (c CoreAuthController) AuthenticateUserWithPassword(CRUD AuthControllerCRUD, username string, password string) (*models.User, common.CustomError) {
	//get the user
	user, err := CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, common.InternalError()
	}

	//check if user was found
	if user == nil {
		return nil, common.ClientError("invalid username and/or password")
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, password)
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return nil, common.ClientError("invalid username and/or password")
	}

	return user, common.NoError()
}
