package router

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/models"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// OAuthAuthenticator is an OAuth implementation of the Authenticator interface.
type OAuthAuthenticator struct{}

// Authenticate attempts to create an access token from the given http request.
func (a OAuthAuthenticator) Authenticate(CRUD models.AccessTokenCRUD, req *http.Request) (*models.AccessToken, requesterror.RequestError) {
	//extract the token string from the authorization header
	splitTokens := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(splitTokens) != 2 {
		return nil, requesterror.ClientError("no bearer token provided")
	}

	//parse the token
	tokenID, err := uuid.Parse(splitTokens[1])
	if err != nil {
		log.Println(common.ChainError("error parsing access token id", err))
		return nil, requesterror.ClientError("bearer token was in invalid format")
	}

	//fetch the token
	token, err := CRUD.GetAccessTokenByID(tokenID)
	if err != nil {
		log.Println(common.ChainError("error getting access token by id", err))
		return nil, requesterror.InternalError()
	}

	// no token found
	if token == nil {
		return nil, requesterror.ClientError("invalid bearer token")
	}

	// auth success
	return token, requesterror.NoError()
}
