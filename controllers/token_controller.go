package controllers

import (
	"authserver/database"
	"authserver/helpers"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TokenController handles requests to "/token" endpoints
type TokenController struct {
	database.UserCRUD
	database.ClientCRUD
	database.ScopeCRUD
	database.AccessTokenCRUD
	helpers.PasswordHasher
}

// PostTokenBody is the struct the body of requests to PostToken should be parsed into
type PostTokenBody struct {
	GrantType string `json:"grant_type"`
	passwordGrantBody
}

// PostToken handles POST requests to "/token"
func (c TokenController) PostToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostTokenBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PostToken request body", err))
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "invalid json body")
		return
	}

	var token *models.AccessToken = nil

	//choose the workflow based on the grant type
	switch body.GrantType {
	case "password":
		token = c.handlePasswordGrant(w, body.passwordGrantBody)
	default:
		sendOAuthErrorResponse(w, http.StatusBadRequest, "unsupported_grant_type", "")
	}

	if token == nil {
		return
	}

	//construct and send the access token response
	sendResponse(w, http.StatusOK, AccessTokenResponse{
		AccessToken: token.ID.String(),
		TokenType:   "bearer",
	})
}

func (c TokenController) handlePasswordGrant(w http.ResponseWriter, body passwordGrantBody) *models.AccessToken {
	//validate parameters
	if body.Username == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing username parameter")
		return nil
	}

	if body.Password == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing password parameter")
		return nil
	}

	if body.ClientID == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing client_id parameter")
		return nil
	}

	if body.Scope == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing scope parameter")
		return nil
	}

	//get the client
	client := parseClient(c.ClientCRUD, w, body.ClientID)
	if client == nil {
		return nil
	}

	//get the scope
	scope := parseScope(c.ScopeCRUD, w, body.Scope)
	if scope == nil {
		return nil
	}

	//get the user
	user, err := c.UserCRUD.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by username", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if user == nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return nil
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, body.Password)
	if err != nil {
		log.Println(helpers.ChainError("error comparing password hashes", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return nil
	}

	//create a new access token
	token := models.CreateNewAccessToken(user.ID, client.ID, scope.ID)

	//save the token
	err = c.AccessTokenCRUD.CreateAccessToken(token)
	if err != nil {
		log.Println(helpers.ChainError("error saving access token", err))
		sendInternalErrorResponse(w)
		return nil
	}

	return token
}
