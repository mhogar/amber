package router

import (
	requesterror "authserver/common/request_error"
	"authserver/models"
	"net/http"
)

// Authenticator is an interface for authenticating and creating an access token from an http request.
type Authenticator interface {
	// Authenticate attempts to create an access token from the given http request.
	Authenticate(CRUD models.AccessTokenCRUD, req *http.Request) (*models.AccessToken, requesterror.RequestError)
}
