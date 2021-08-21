package controllers

import (
	"authserver/common"
	"authserver/models"
	"log"
)

type CoreSessionController struct {
	AuthController AuthController
}

func (c CoreSessionController) CreateSession(CRUD SessionControllerCRUD, username string, password string) (*models.Session, common.CustomError) {
	//authenticate the user
	user, cerr := c.AuthController.AuthenticateUserWithPassword(CRUD, username, password)
	if cerr.Type != common.ErrorTypeNone {
		return nil, cerr
	}

	//create a new session
	session := models.CreateNewSession(user)

	//save the session
	err := CRUD.SaveSession(session)
	if err != nil {
		log.Println(common.ChainError("error saving session", err))
		return nil, common.InternalError()
	}

	return session, common.NoError()
}

func (c CoreSessionController) DeleteSession(CRUD SessionControllerCRUD, session *models.Session) common.CustomError {
	//delete the session
	err := CRUD.DeleteSession(session)
	if err != nil {
		log.Println(common.ChainError("error deleting session", err))
		return common.InternalError()
	}

	return common.NoError()
}

func (c CoreSessionController) DeleteAllOtherUserSessions(CRUD SessionControllerCRUD, session *models.Session) common.CustomError {
	//delete the session
	err := CRUD.DeleteAllOtherUserSessions(session)
	if err != nil {
		log.Println(common.ChainError("error deleting all other user sessions", err))
		return common.InternalError()
	}

	return common.NoError()
}
