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

type ClientDataResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostClientBody struct {
	Name string `json:"name"`
}

func (h CoreHandlers) PostClient(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PostClientBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostClient request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the client
	client, cerr := h.Controllers.CreateClient(CRUD, body.Name)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return newClientDataResponse(client)
}

type PutClientBody struct {
	Name string `json:"name"`
}

func (h CoreHandlers) PutClient(req *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PutClientBody

	//parse the id
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
	cerr := h.Controllers.UpdateClient(CRUD, client)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return newClientDataResponse(client)
}

func (h CoreHandlers) DeleteClient(_ *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//parse the id
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//delete the client
	cerr := h.Controllers.DeleteClient(CRUD, id)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func newClientDataResponse(client *models.Client) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(ClientDataResponse{
		ID:   client.UID.String(),
		Name: client.Name,
	})
}
