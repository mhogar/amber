package handlers

import (
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/parsers"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ClientDataResponse struct {
	ID string `json:"id"`
	PostClientBody
}

func (h CoreAPIHandlers) GetClients(_ *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//get the clients
	clients, cerr := h.Controllers.GetClients(CRUD)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//return the data
	data := make([]ClientDataResponse, len(clients))
	for index, client := range clients {
		data[index] = h.newClientDataResponse(client)
	}
	return common.NewSuccessDataResponse(data)
}

type PostClientBody struct {
	Name        string `json:"name"`
	RedirectUrl string `json:"redirect_url"`
	TokenType   int    `json:"token_type"`
	KeyUri      string `json:"key_uri"`
}

func (h CoreAPIHandlers) PostClient(req *http.Request, _ httprouter.Params, _ *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	var body PostClientBody

	//parse the body
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostClient request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//create the client model
	client := models.CreateNewClient(body.Name, body.RedirectUrl, body.TokenType, body.KeyUri)

	//create the client
	cerr := h.Controllers.CreateClient(CRUD, client)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newClientDataResponse(client))
}

func (h CoreAPIHandlers) PutClient(req *http.Request, params httprouter.Params, _ *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	var body PostClientBody

	//parse the id
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		log.Println(common.ChainError("error parsing id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//parse the body
	err = parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PutClient request body", err))
		return common.NewBadRequestResponse("invalid request body")
	}

	//create the client model
	client := models.CreateClient(id, body.Name, body.RedirectUrl, body.TokenType, body.KeyUri)

	//update the client
	cerr := h.Controllers.UpdateClient(CRUD, client)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessDataResponse(h.newClientDataResponse(client))
}

func (h CoreAPIHandlers) DeleteClient(_ *http.Request, params httprouter.Params, _ *models.Session, _ parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
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

func (CoreAPIHandlers) newClientDataResponse(client *models.Client) ClientDataResponse {
	return ClientDataResponse{
		ID: client.UID.String(),
		PostClientBody: PostClientBody{
			Name:        client.Name,
			RedirectUrl: client.RedirectUrl,
			TokenType:   client.TokenType,
			KeyUri:      client.KeyUri,
		},
	}
}
