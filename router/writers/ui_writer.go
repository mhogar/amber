package writers

import (
	"net/http"

	"github.com/mhogar/amber/common"
)

type UIResponseWriter struct{}

func NewUIResponseWriter() UIResponseWriter {
	return UIResponseWriter{}
}

func (UIResponseWriter) WriteResponse(w http.ResponseWriter, status int, res interface{}) {
	w.WriteHeader(status)
	w.Write(res.([]byte))
}

func (rw UIResponseWriter) WriteErrorResponse(w http.ResponseWriter, _ int, message string) {
	//TODO: render generic error view
}

func (rw UIResponseWriter) WriteInternalErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInternalServerErrorResponse()
	rw.WriteErrorResponse(w, status, res.Error)
}

func (rw UIResponseWriter) WriteInsufficientPermissionsErrorResponse(w http.ResponseWriter) {
	status, res := common.NewInsufficientPermissionsErrorResponse()
	rw.WriteErrorResponse(w, status, res.Error)
}
