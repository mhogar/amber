package handlers

import (
	"authserver/common"
	"authserver/config"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type TokenViewData struct {
	AppName  string
	ClientID string
	Error    string
}

func (h CoreHandlers) GetToken(req *http.Request, _ httprouter.Params, _ *models.Session, _ data.DataCRUD) (int, interface{}) {
	return h.renderTokenView(req.URL.Query().Get("client_id"), "")
}

func (h CoreHandlers) PostToken(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//get the form values
	clientIdStr := req.PostFormValue("client_id")
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	//parse the client id
	clientID, err := uuid.Parse(clientIdStr)
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return h.renderTokenView(clientIdStr, "client_id is not provided or in an invalid format")
	}

	//create the token redirect url
	redirectUrl, cerr := h.Controllers.CreateTokenRedirectURL(CRUD, clientID, username, password)
	if cerr.Type != common.ErrorTypeNone {
		return h.renderTokenView(clientIdStr, cerr.Error())
	}

	//send redirect response
	return http.StatusSeeOther, redirectUrl
}

func (h CoreHandlers) renderTokenView(clientID string, errMessage string) (int, interface{}) {
	//fill in the data struct
	data := TokenViewData{
		AppName:  config.GetAppName(),
		ClientID: clientID,
		Error:    errMessage,
	}

	//render the view
	return http.StatusOK, h.Renderer.RenderView("token.gohtml", data)
}
