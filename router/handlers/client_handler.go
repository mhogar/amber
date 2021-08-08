package handlers

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type PostClientBody struct {
	Name string `json:"name"`
}

func (h CoreHandlers) PostClient(req *http.Request, _ httprouter.Params, _ *models.AccessToken, tx data.Transaction) (int, interface{}) {
	var body PostClientBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostClient request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the client
	_, rerr := h.Controllers.CreateClient(tx, body.Name)
	if rerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

type PutClientBody struct {
	Name string `json:"name"`
}

func (h CoreHandlers) PutClient(req *http.Request, params httprouter.Params, _ *models.AccessToken, tx data.Transaction) (int, interface{}) {
	var body PutClientBody

	// parse the id
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//parse the body
	err = parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PutClient request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the client model
	client := models.CreateClient(id, body.Name)

	//update the client
	rerr := h.Controllers.UpdateClient(tx, client)
	if rerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (h CoreHandlers) DeleteClient(_ *http.Request, params httprouter.Params, _ *models.AccessToken, tx data.Transaction) (int, interface{}) {
	// parse the id
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//delete the client
	rerr := h.Controllers.DeleteClient(tx, id)
	if rerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(rerr.Error())
	}
	if rerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}