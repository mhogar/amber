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
	ID string `json:"id"`
	PostClientBody
}

type PostClientBody struct {
	Name        string `json:"name"`
	RedirectUrl string `json:"redirect_url"`
	TokenType   int    `json:"token_type"`
	KeyUri      string `json:"key_uri"`
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
	client, cerr := h.Controllers.CreateClient(CRUD, body.Name, body.RedirectUrl, body.TokenType, body.KeyUri)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return h.newClientDataResponse(client)
}

func (h CoreHandlers) PutClient(req *http.Request, params httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PostClientBody

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
	client := models.CreateClient(id, body.Name, body.RedirectUrl, body.TokenType, body.KeyUri)

	//update the client
	cerr := h.Controllers.UpdateClient(CRUD, client)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return h.newClientDataResponse(client)
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

func (CoreHandlers) newClientDataResponse(client *models.Client) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(ClientDataResponse{
		ID: client.UID.String(),
		PostClientBody: PostClientBody{
			Name:        client.Name,
			RedirectUrl: client.RedirectUrl,
			TokenType:   client.TokenType,
			KeyUri:      client.KeyUri,
		},
	})
}
