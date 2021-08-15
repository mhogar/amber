package common

import "net/http"

// BasicResponse represents a response with a simple true/false success field.
type BasicResponse struct {
	Success bool `json:"success"`
}

func NewSuccessResponse() (int, BasicResponse) {
	return http.StatusOK, BasicResponse{
		Success: true,
	}
}

// ErrorResponse represents a response with a true/false success field and an error message.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
	}
}

func NewBadRequestResponse(err string) (int, ErrorResponse) {
	return http.StatusBadRequest, NewErrorResponse(err)
}

func NewInternalServerErrorResponse() (int, ErrorResponse) {
	return http.StatusInternalServerError, NewErrorResponse("an internal error occurred")
}

// DataResponse represents a response with a true/false success field and generic data.
type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func NewSuccessDataResponse(data interface{}) (int, DataResponse) {
	return http.StatusOK, DataResponse{
		Success: true,
		Data:    data,
	}
}

// OAuthErrorResponse represents an error response defined by the oauth spec.
type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewOAuthErrorResponse(name string, description string) (int, OAuthErrorResponse) {
	return http.StatusBadRequest, OAuthErrorResponse{
		Error:            name,
		ErrorDescription: description,
	}
}

// AccessTokenResponse represents an access token response defined by the oauth spec.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewAccessTokenResponse(token string) (int, AccessTokenResponse) {
	return http.StatusOK, AccessTokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
	}
}
