package writers

import (
	"encoding/json"
	"net/http"

	"github.com/mhogar/amber/common"
)

type JSONResponseWriter struct{}

func NewJSONResponseWriter() JSONResponseWriter {
	return JSONResponseWriter{}
}

func (JSONResponseWriter) WriteResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	if res != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}
}

func (rw JSONResponseWriter) WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	rw.WriteResponse(w, status, common.NewErrorResponse(message))
}

func (rw JSONResponseWriter) WriteInternalErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInternalServerErrorResponse()
	rw.WriteErrorResponse(w, status, res.Error)
}

func (rw JSONResponseWriter) WriteInsufficientPermissionsErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInsufficientPermissionsErrorResponse()
	rw.WriteErrorResponse(w, status, res.Error)
}
