package writers

import "net/http"

type ResponseWriter interface {
	WriteResponse(w http.ResponseWriter, status int, res interface{})
	WriteErrorResponse(w http.ResponseWriter, status int, message string)
	WriteInternalErrorResponse(w http.ResponseWriter)
	WriteInsufficientPermissionsErrorResponse(w http.ResponseWriter)
}
