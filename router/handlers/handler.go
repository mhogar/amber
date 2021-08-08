package handlers

import (
	"authserver/controllers"
	"authserver/data"
	"authserver/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handlers interface {
	// PostUser handles POST requests to /user
	PostUser(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// DeleteUser handles DELETE requests to /user
	DeleteUser(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// PatchUserPassword handles PATCH requests to /user/password
	PatchUserPassword(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// PostToken handles POST requests to /token
	PostToken(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// DeleteToken handles DELETE requests to /token
	DeleteToken(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// PostClient handles POST requests to /client
	PostClient(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// PutClient handles PUT requests to /client
	PutClient(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

	// DeleteClient handles DELETE requests to /client
	DeleteClient(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})
}

type CoreHandlers struct {
	Controllers controllers.Controllers
}
