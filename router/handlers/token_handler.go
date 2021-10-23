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

type TokenViewData struct {
	ClientID string
	Error    string
}

func (h CoreUIHandlers) GetToken(req *http.Request, _ httprouter.Params, _ *models.Session, _ parsers.BodyParser, _ data.DataCRUD) (int, interface{}) {
	return h.renderTokenView(req, req.URL.Query().Get("client_id"), "")
}

func (h CoreUIHandlers) PostToken(req *http.Request, _ httprouter.Params, _ *models.Session, _ parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//get the form values
	clientIdStr := req.PostFormValue("client_id")
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	//parse the client id
	clientID, err := uuid.Parse(clientIdStr)
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return h.renderTokenView(req, clientIdStr, "client_id is not provided or in an invalid format")
	}

	//create the token redirect url
	redirectUrl, cerr := h.Controllers.CreateTokenRedirectURL(CRUD, clientID, username, password)
	if cerr.Type != common.ErrorTypeNone {
		return h.renderTokenView(req, clientIdStr, cerr.Error())
	}

	//send redirect response
	return common.NewRedirectResponse(redirectUrl, "")
}

func (h CoreUIHandlers) renderTokenView(req *http.Request, clientID string, errMessage string) (int, interface{}) {
	//fill in the data struct
	data := TokenViewData{
		ClientID: clientID,
		Error:    errMessage,
	}

	//render the view
	return http.StatusOK, h.Renderer.RenderView(req, data, "token/index", "partials/login_form", "partials/alert")
}
