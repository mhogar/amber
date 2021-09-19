package router

import (
	"encoding/json"
	"net/http"

	"github.com/mhogar/amber/common"
)

func sendRawResponse(w http.ResponseWriter, status int, res []byte) {
	w.WriteHeader(status)
	w.Write(res)
}

func sendJSONResponse(w http.ResponseWriter, status int, res interface{}) {
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
	sendJSONResponse(w, status, common.NewErrorResponse(message))
}

func sendInternalErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInternalServerErrorResponse()
	sendErrorResponse(w, status, res.Error)
}

func sendInsufficientPermissionsErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInsufficientPermissionsErrorResponse()
	sendErrorResponse(w, status, res.Error)
}
