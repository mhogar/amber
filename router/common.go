package router

import (
	"log"
	"net/http"
)

func panicHandler(w http.ResponseWriter, req *http.Request, info interface{}) {
	log.Println(info)
	sendInternalErrorResponse(w)
}
