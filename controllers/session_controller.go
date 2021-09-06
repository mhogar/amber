package controllers

import (
	"authserver/common"
	"authserver/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CoreSessionController struct {
	AuthController AuthController
}

func (c CoreSessionController) CreateSession(CRUD SessionControllerCRUD, username string, password string) (*models.Session, common.CustomError) {
	//authenticate the user
	_, cerr := c.AuthController.AuthenticateUserWithPassword(CRUD, username, password)
	if cerr.Type != common.ErrorTypeNone {
		return nil, cerr
	}

	//create a new session
	session := models.CreateNewSession(username, 0)

	//save the session
	err := CRUD.SaveSession(session)
	if err != nil {
		log.Println(common.ChainError("error saving session", err))
		return nil, common.InternalError()
	}

	return session, common.NoError()
}

func (c CoreSessionController) DeleteSession(CRUD SessionControllerCRUD, id uuid.UUID) common.CustomError {
	//delete the session
	res, err := CRUD.DeleteSession(id)
	if err != nil {
		log.Println(common.ChainError("error deleting session", err))
		return common.InternalError()
	}

	//verify sesion was actually found
	if !res {
		return common.ClientError(fmt.Sprintf("session with id %s not found", id.String()))
	}

	return common.NoError()
}

func (c CoreSessionController) DeleteAllOtherUserSessions(CRUD SessionControllerCRUD, username string, id uuid.UUID) common.CustomError {
	//delete the session
	err := CRUD.DeleteAllOtherUserSessions(username, id)
	if err != nil {
		log.Println(common.ChainError("error deleting all other user sessions", err))
		return common.InternalError()
	}

	return common.NoError()
}
