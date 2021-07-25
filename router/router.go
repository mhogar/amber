package router

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

type IRouterFactory interface {
	CreateRouter() *httprouter.Router
}

type RouterFactory struct {
	Authenticator Authenticator
	ScopeFactory  data.IScopeFactory
	Handlers      IHandlers
}

// CreateRouter creates a new httprouter with the endpoints and panic handler configured.
func (rf RouterFactory) CreateRouter() *httprouter.Router {
	r := httprouter.New()
	r.PanicHandler = panicHandler

	//user routes
	r.POST("/user", rf.createHandler(rf.Handlers.PostUser, false))
	r.DELETE("/user", rf.createHandler(rf.Handlers.DeleteUser, true))
	r.PATCH("/user/password", rf.createHandler(rf.Handlers.PatchUserPassword, true))

	//token routes
	r.POST("/token", rf.createHandler(rf.Handlers.PostToken, false))
	r.DELETE("/token", rf.createHandler(rf.Handlers.DeleteToken, true))

	return r
}
