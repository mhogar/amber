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

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	sendResponse(w, status, common.NewErrorResponse(message))
}

func sendInternalErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInternalServerErrorResponse()
	sendErrorResponse(w, status, res.Error)
}

func sendInsufficientPermissionsErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInsufficientPermissionsErrorResponse()
	sendErrorResponse(w, status, res.Error)
}
