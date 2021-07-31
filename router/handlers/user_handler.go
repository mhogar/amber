package handlers

import (
	"log"
	"net/http"

	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/data"
	"authserver/models"

	"github.com/julienschmidt/httprouter"
)

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (h Handlers) PostUser(req *http.Request, _ httprouter.Params, _ *models.AccessToken, tx data.Transaction) (int, interface{}) {
	//parse the body
	var body PostUserBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the user
	_, rerr := h.Controllers.CreateUser(tx, body.Username, body.Password)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

// DeleteUser handles DELETE requests to "/user"
func (h Handlers) DeleteUser(_ *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	//delete the user
	rerr := h.Controllers.DeleteUser(tx, token.User)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

// PatchUserPasswordBody is the struct the body of requests to PatchUserPassword should be parsed into
type PatchUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// PatchUserPassword handles PATCH requests to "/user/password"
func (h Handlers) PatchUserPassword(req *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	//parse the body
	var body PatchUserPasswordBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchUserPassword request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//update the password
	rerr := h.Controllers.UpdateUserPassword(tx, token.User, body.OldPassword, body.NewPassword)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//delete all other user access tokens
	rerr = h.Controllers.DeleteAllOtherUserTokens(tx, token)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}
