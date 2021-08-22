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
