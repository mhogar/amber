package handlers

import (
	"log"
	"net/http"

	"authserver/common"
	"authserver/data"
	"authserver/models"

	"github.com/julienschmidt/httprouter"
)

type UserDataResponse struct {
	Username string `json:"username"`
	PutUserBody
}

type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rank     int    `json:"rank"`
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
	user, cerr := h.Controllers.CreateUser(CRUD, body.Username, body.Password, body.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return h.newUserDataResponse(user)
}

type PutUserBody struct {
	Rank int `json:"rank"`
}

func (h CoreHandlers) PutUser(req *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//parse the body
	var body PutUserBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PutUser request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//update the user
	user, cerr := h.Controllers.UpdateUser(CRUD, username, body.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return h.newUserDataResponse(user)
}

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

func (h CoreHandlers) DeleteUser(_ *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//delete the user
	cerr := h.Controllers.DeleteUser(CRUD, username)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (CoreHandlers) newUserDataResponse(user *models.User) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(UserDataResponse{
		Username: user.Username,
		PutUserBody: PutUserBody{
			Rank: user.Rank,
		},
	})
}
