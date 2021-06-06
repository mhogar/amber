package router

import (
	"authserver/common"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func parseJSONBody(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		log.Println(err)
		return errors.New("invalid request body")
	}

	return nil
}

func sendResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Panic(err) //panic if can't write response
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
