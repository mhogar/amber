package handlers

import (
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/parsers"

	"github.com/julienschmidt/httprouter"
)

type UserDataResponse struct {
	Username string `json:"username"`
	PutUserBody
}

func (h CoreAPIHandlers) GetUsers(_ *http.Request, _ httprouter.Params, session *models.Session, _ parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//get the users
	users, cerr := h.Controllers.GetUsersWithLesserRank(CRUD, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//return the data
	data := make([]UserDataResponse, len(users))
	for index, user := range users {
		data[index] = h.newUserDataResponse(user)
	}
	return common.NewSuccessDataResponse(data)
}

type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rank     int    `json:"rank"`
}

func (h CoreAPIHandlers) PostUser(req *http.Request, _ httprouter.Params, session *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//parse the body
	var body PostUserBody
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//verify the session has a greater rank the user being created
	if body.Rank > session.Rank {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//create the user
	user, cerr := h.Controllers.CreateUser(CRUD, body.Username, body.Password, body.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newUserDataResponse(user))
}

type PutUserBody struct {
	Rank int `json:"rank"`
}

func (h CoreAPIHandlers) PutUser(req *http.Request, params httprouter.Params, session *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//parse the body
	var body PutUserBody
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PutUser request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//verify the session has a greater rank than the new rank
	if body.Rank > session.Rank {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//verify the session has a greater rank than the user's current rank
	res, cerr := h.Controllers.VerifyUserRank(CRUD, username, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}
	if !res {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//update the user
	user, cerr := h.Controllers.UpdateUser(CRUD, username, body.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newUserDataResponse(user))
}

type PatchPasswordBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h CoreAPIHandlers) PatchPassword(req *http.Request, _ httprouter.Params, session *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//parse the body
	var body PatchPasswordBody
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchPassword request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//update the password
	cerr := h.Controllers.UpdateUserPasswordWithAuth(CRUD, session.Username, body.OldPassword, body.NewPassword)
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

type PatchUserPasswordBody struct {
	Password string `json:"password"`
}

func (h CoreAPIHandlers) PatchUserPassword(req *http.Request, params httprouter.Params, session *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//parse the body
	var body PatchUserPasswordBody
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchUserPasswordBody request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//verify the session has a greater rank than the user
	res, cerr := h.Controllers.VerifyUserRank(CRUD, username, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}
	if !res {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//update the password
	cerr = h.Controllers.UpdateUserPassword(CRUD, username, body.Password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//delete all user sessions
	cerr = h.Controllers.DeleteAllUserSessions(CRUD, username)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (h CoreAPIHandlers) DeleteUser(_ *http.Request, params httprouter.Params, session *models.Session, _ parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//verify the session has a greater rank than the user
	res, cerr := h.Controllers.VerifyUserRank(CRUD, username, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}
	if !res {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//delete the user
	cerr = h.Controllers.DeleteUser(CRUD, username)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (CoreAPIHandlers) newUserDataResponse(user *models.User) UserDataResponse {
	return UserDataResponse{
		Username: user.Username,
		PutUserBody: PutUserBody{
			Rank: user.Rank,
		},
	}
}
