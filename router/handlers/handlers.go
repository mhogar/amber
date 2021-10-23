package handlers

import (
	"net/http"

	"github.com/mhogar/amber/controllers"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/parsers"
	"github.com/mhogar/amber/router/renderer"

	"github.com/julienschmidt/httprouter"
)

type APIHandlers interface {
	// GetUsers handles GET requests to /api/users.
	GetUsers(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostUser handles POST requests to /api/user.
	PostUser(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PutUser handles PUT requests to /api/user/:username.
	PutUser(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PatchPassword handles PATCH requests to /api/user/password.
	PatchPassword(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PatchUserPassword handles PATCH requests to /api/user/:username/password.
	PatchUserPassword(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// DeleteUser handles DELETE requests to /api/user/:username.
	DeleteUser(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// GetClients handles GET requests to /api/clients.
	GetClients(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostClient handles POST requests to /api/client.
	PostClient(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PutClient handles PUT requests to /api/client/:id.
	PutClient(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// DeleteClient handles DELETE requests to /api/client/:id.
	DeleteClient(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// GetUserRoles handles GET requests to /api/client/:id/roles.
	GetUserRoles(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostUserRole handles POST requests to /api/client/:id/role.
	PostUserRole(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PutUserRole handles PUT requests to /api/client/:id/role/:username.
	PutUserRole(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// DeleteUserRole handles DELETE requests to /api/client/:id/role/:username.
	DeleteUserRole(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostSession handles POST requests to /api/session.
	PostSession(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// DeleteSession handles DELETE requests to /api/session.
	DeleteSession(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})
}

type UIHandlers interface {
	// GetHome Handles GET requests to /.
	GetHome(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// GetLogin Handles GET requests to /login.
	GetLogin(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostSession Handles POST requests to /session.
	PostSession(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// GetToken handles GET requests to /token.
	GetToken(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

	// PostToken handles POST requests to /token.
	PostToken(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})
}

type CoreAPIHandlers struct {
	Controllers controllers.Controllers
}

type CoreUIHandlers struct {
	API         APIHandlers
	Controllers controllers.Controllers
	Renderer    renderer.Renderer
}
