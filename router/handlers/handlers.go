package handlers

import (
	"authserver/controllers"
	"authserver/data"
	"authserver/models"
	"authserver/router/renderer"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handlers interface {
	// GetUsers handles GET requests to /users.
	GetUsers(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostUser handles POST requests to /user.
	PostUser(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PutUser handles PUT requests to /user/:username.
	PutUser(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PatchPassword handles PATCH requests to /user/password.
	PatchPassword(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PatchUserPassword handles PATCH requests to /user/:username/password.
	PatchUserPassword(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteUser handles DELETE requests to /user/:username.
	DeleteUser(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// GetClients handles GET requests to /clients.
	GetClients(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostClient handles POST requests to /client.
	PostClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PutClient handles PUT requests to /client/:id.
	PutClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteClient handles DELETE requests to /client/:id.
	DeleteClient(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// GetUserRoles handles GET requests to /user/:username/roles.
	GetUserRoles(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostUserRole handles POST requests to /user/:username/role.
	PostUserRole(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PutUserRole handles PUT requests to /user/:username/role/:client_id.
	PutUserRole(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteUserRole handles DELETE requests to /user/:username/role/:client_id.
	DeleteUserRole(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostSession handles POST requests to /session.
	PostSession(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// DeleteSession handles DELETE requests to /session.
	DeleteSession(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// GetToken handles GET requests to /token.
	GetToken(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

	// PostToken handles POST requests to /token.
	PostToken(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})
}

type CoreHandlers struct {
	Controllers controllers.Controllers
	Renderer    renderer.Renderer
}
