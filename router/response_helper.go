package router

import (
	"authserver/common"
	"encoding/json"
	"net/http"
)

func sendResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	if res != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}
}

func sendErrorResponse(w http.ResponseWriter, status int, messsage string) {
	sendResponse(w, status, common.ErrorResponse{
		Success: false,
		Error:   messsage,
	})
}

func sendInternalErrorResponse(w http.ResponseWriter) {
	sendErrorResponse(w, http.StatusInternalServerError, "an internal error occurred")
}

func sendInsufficientPermissionsErrorResponse(w http.ResponseWriter) {
	sendErrorResponse(w, http.StatusForbidden, "insufficient permissions to perform the requested action")
}
