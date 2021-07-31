package handlers

import (
	"authserver/controllers"
	"authserver/data"
	"authserver/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IHandlers interface {
	PostUser(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})
	DeleteUser(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})
	PatchUserPassword(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	PostToken(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})
	DeleteToken(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})
}

type Handlers struct {
	Controllers controllers.Controllers
}
