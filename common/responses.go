package common

import "net/http"

// BasicResponse represents a response with a simple true/false success field.
type BasicResponse struct {
	Success bool `json:"success"`
}

// NewSuccessResponse returns an http OK status and a success basic response.
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

// NewErrorResponse creates a new ErrorResponse with a false success field and the provided error message.
func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
	}
}

// NewBadRequestResponse returns an http BadRequest status and a new error response with the provided error message.
func NewBadRequestResponse(err string) (int, ErrorResponse) {
	return http.StatusBadRequest, NewErrorResponse(err)
}

// NewInternalServerErrorResponse returns an http StatusInternalServerError status and a new error response with an internal error message.
func NewInternalServerErrorResponse() (int, ErrorResponse) {
	return http.StatusInternalServerError, NewErrorResponse(InternalError().Error())
}

// NewInsufficientPermissionsErrorResponse returns an http StatusForbidden status and a new error response with an insufficient permissions error message.
func NewInsufficientPermissionsErrorResponse() (int, ErrorResponse) {
	return http.StatusForbidden, NewErrorResponse("insufficient permissions to perform the requested action")
}

// DataResponse represents a response with a true/false success field and generic data.
type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// NewSuccessDataResponse return an http OK status and a new DataResponse with the provided data.
func NewSuccessDataResponse(data interface{}) (int, DataResponse) {
	return http.StatusOK, DataResponse{
		Success: true,
		Data:    data,
	}
}
