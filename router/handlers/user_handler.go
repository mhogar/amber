package handlers

import (
	"log"
	"net/http"

	"authserver/common"
	"authserver/data"
	"authserver/models"

	"github.com/julienschmidt/httprouter"
)

// PostUserBody is the struct the body of requests to PostUser should be parsed into.
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h CoreHandlers) PostUser(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//parse the body
	var body PostUserBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the user
	_, cerr := h.Controllers.CreateUser(CRUD, body.Username, body.Password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (h CoreHandlers) DeleteUser(_ *http.Request, _ httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//delete the user
	cerr := h.Controllers.DeleteUser(CRUD, session.Username)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

// PatchUserPasswordBody is the struct the body of requests to PatchUserPassword should be parsed into.
type PatchUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h CoreHandlers) PatchUserPassword(req *http.Request, _ httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//parse the body
	var body PatchUserPasswordBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchUserPassword request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//update the password
	cerr := h.Controllers.UpdateUserPassword(CRUD, session.Username, body.OldPassword, body.NewPassword)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//delete all other user sessions
	cerr = h.Controllers.DeleteAllOtherUserSessions(CRUD, session.Username, session.Token)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}
