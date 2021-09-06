package handlers

import (
	"authserver/controllers"
	"authserver/data"
	"authserver/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handlers interface {
	// PostUser handles POST requests to /user.
	PostUser(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteUser handles DELETE requests to /user.
	DeleteUser(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PatchUserPassword handles PATCH requests to /user/password.
	PatchUserPassword(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostClient handles POST requests to /client.
	PostClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PutClient handles PUT requests to /client.
	PutClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteClient handles DELETE requests to /client.
	DeleteClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PutClientRoles handles PUT requests to /client/:id/roles.
	PutClientRoles(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostSession handles POST requests to /token.
	PostSession(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteSession handles DELETE requests to /token.
	DeleteSession(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostToken handles POST requests to /token.
	PostToken(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})
}

type CoreHandlers struct {
	Controllers controllers.Controllers
}
