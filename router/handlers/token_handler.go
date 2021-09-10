package handlers

import (
	"authserver/common"
	"authserver/config"
	"authserver/data"
	"authserver/models"
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type TokenViewData struct {
	AppName  string
	ClientID string
}

func (h CoreHandlers) GetToken(req *http.Request, _ httprouter.Params, _ *models.Session, _ data.DataCRUD) (int, interface{}) {
	//parse the template
	t := template.Must(template.ParseFiles(path.Join(config.GetAppRoot(), "views", "token.gohtml")))

	//fill in the data struct
	data := TokenViewData{
		AppName:  config.GetAppName(),
		ClientID: req.URL.Query().Get("client_id"),
	}

	//render the template
	var buffer bytes.Buffer
	err := t.Execute(&buffer, data)
	if err != nil {
		log.Println("error rendering get token template")
		return common.NewInternalServerErrorResponse()
	}

	return http.StatusOK, buffer.Bytes()
}

type PostTokenBody struct {
	ClientId uuid.UUID `json:"client_id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func (h CoreHandlers) PostToken(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//get the form values
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	//parse the client id
	clientID, err := uuid.Parse(req.PostFormValue("client_id"))
	if err != nil {
		log.Println(common.ChainError("error parsing client id", err))
		return common.NewBadRequestResponse("client id is in an invalid format")
	}

	//create the token redirect url
	redirectUrl, cerr := h.Controllers.CreateTokenRedirectURL(CRUD, clientID, username, password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	//send redirect response
	return http.StatusSeeOther, redirectUrl
}
