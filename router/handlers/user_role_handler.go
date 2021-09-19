package handlers

import (
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserRoleDataResponse struct {
	PostUserRoleBody
}

func (h CoreHandlers) GetUserRoles(_ *http.Request, params httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//parse the client id
	clientID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//get the roles
	roles, cerr := h.Controllers.GetUserRolesWithLesserRankByClientUID(CRUD, clientID, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//return the data
	data := make([]UserRoleDataResponse, len(roles))
	for index, role := range roles {
		data[index] = h.newUserRoleDataResponse(role)
	}
	return common.NewSuccessDataResponse(data)
}

type PostUserRoleBody struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h CoreHandlers) PostUserRole(req *http.Request, params httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PostUserRoleBody

	//parse the client id
	clientID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//parse the body
	err = parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUserRoleBody request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//verify the session has a greater rank than the user
	res, cerr := h.Controllers.VerifyUserRank(CRUD, body.Username, session.Rank)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}
	if !res {
		return common.NewInsufficientPermissionsErrorResponse()
	}

	//create the model
	role := models.CreateUserRole(clientID, body.Username, body.Role)

	//create the user-role
	cerr = h.Controllers.CreateUserRole(CRUD, role)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newUserRoleDataResponse(role))
}

type PutUserRoleBody struct {
	Role string `json:"role"`
}

func (h CoreHandlers) PutUserRole(req *http.Request, params httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PutUserRoleBody

	//parse the client id
	clientID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//get the username
	username := params.ByName("username")
	if username == "" {
		return common.NewBadRequestResponse("username not provided")
	}

	//parse the body
	err = parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PutUserRoleBody request body", err))
		return common.NewBadRequestResponse("invalid json body")
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

	//create the model
	role := models.CreateUserRole(clientID, username, body.Role)

	//update the user-role
	cerr = h.Controllers.UpdateUserRole(CRUD, role)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newUserRoleDataResponse(role))
}

func (h CoreHandlers) DeleteUserRole(_ *http.Request, params httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//parse the client id
	clientID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

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

	//delete the user-role
	cerr = h.Controllers.DeleteUserRole(CRUD, username, clientID)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (CoreHandlers) newUserRoleDataResponse(role *models.UserRole) UserRoleDataResponse {
	return UserRoleDataResponse{
		PostUserRoleBody: PostUserRoleBody{
			Username: role.Username,
			Role:     role.Role,
		},
	}
}
