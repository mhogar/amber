package handlers

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// PostTokenBody is the struct the body of requests to PostToken should be parsed into
type PostTokenBody struct {
	GrantType string `json:"grant_type"`
	PostTokenPasswordGrantBody
}

// PostTokenPasswordGrantBody is the struct the body of password grant requests to PostToken should be parsed into
type PostTokenPasswordGrantBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

// PostToken handles POST requests to "/token"
func (h Handlers) PostToken(req *http.Request, _ httprouter.Params, _ *models.AccessToken, tx data.Transaction) (int, interface{}) {
	var body PostTokenBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostToken request body", err))
		return common.NewOAuthErrorResponse("invalid_request", "invalid json body")
	}

	//validate grant type is present
	if body.GrantType == "" {
		return common.NewOAuthErrorResponse("invalid_request", "missing grant_type parameter")
	}

	//choose the workflow based on the grant type
	switch body.GrantType {
	case "password":
		return h.handlePasswordGrant(body.PostTokenPasswordGrantBody, tx)
	default:
		return common.NewOAuthErrorResponse("unsupported_grant_type", "")
	}
}

func (h Handlers) handlePasswordGrant(body PostTokenPasswordGrantBody, tx data.Transaction) (int, interface{}) {
	//validate parameters
	if body.Username == "" {
		return common.NewOAuthErrorResponse("invalid_request", "missing username parameter")
	}
	if body.Password == "" {
		return common.NewOAuthErrorResponse("invalid_request", "missing password parameter")
	}
	if body.ClientID == "" {
		return common.NewOAuthErrorResponse("invalid_request", "missing client_id parameter")
	}
	if body.Scope == "" {
		return common.NewOAuthErrorResponse("invalid_request", "missing scope parameter")
	}

	//parse the client id
	clientID, err := uuid.Parse(body.ClientID)
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewOAuthErrorResponse("invalid_client", "client_id was in invalid format")
	}

	//create the token
	token, rerr := h.Controllers.CreateTokenFromPassword(tx, body.Username, body.Password, clientID, body.Scope)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewOAuthErrorResponse(rerr.ErrorName, rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewAccessTokenResponse(token.ID.String())
}

// DeleteToken handles DELETE requests to "/token"
func (h Handlers) DeleteToken(_ *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	//delete the token
	rerr := h.Controllers.DeleteToken(tx, token)
	if rerr.Type == requesterror.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == requesterror.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}
