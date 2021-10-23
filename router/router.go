package router

import (
	"log"
	"net/http"
	"strings"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/handlers"
	"github.com/mhogar/amber/router/parsers"
	"github.com/mhogar/amber/router/writers"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type RouterFactory interface {
	// CreateRouter creates a new httprouter with the endpoints and panic handler configured.
	CreateRouter() *httprouter.Router
}

type CoreRouterFactory struct {
	ScopeFactory data.ScopeFactory
	UIHandlers   handlers.UIHandlers
	APIHandlers  handlers.APIHandlers
}

func (rf CoreRouterFactory) CreateRouter() *httprouter.Router {
	r := httprouter.New()

	//host public folder as file server
	r.ServeFiles("/public/*filepath", http.Dir(config.GetAppRoot("public")))

	//token routes
	r.GET("/token", rf.createUIHandler(rf.UIHandlers.GetToken, false, 0))
	r.POST("/token", rf.createUIHandler(rf.UIHandlers.PostToken, false, 0))

	//other routes
	rf.createUIRoutes(r)
	rf.createAPIRoutes(r)

	//set panic handler
	rf.setPanicHandler(r, writers.NewJSONResponseWriter())

	return r
}

func (rf CoreRouterFactory) createUIRoutes(r *httprouter.Router) {
	//home routes
	r.GET("/", rf.createUIHandler(rf.UIHandlers.GetHome, false, 0))

	//login routes
	r.GET("/login", rf.createUIHandler(rf.UIHandlers.GetLogin, false, 0))
}

func (rf CoreRouterFactory) createAPIRoutes(r *httprouter.Router) {
	//user routes
	r.GET("/api/users", rf.createAPIHandler(rf.APIHandlers.GetUsers, true, 0))
	r.POST("/api/user", rf.createAPIHandler(rf.APIHandlers.PostUser, true, 0))
	r.PUT("/api/user/:username", rf.createAPIHandler(rf.APIHandlers.PutUser, true, 0))
	r.PATCH("/api/user/password", rf.createAPIHandler(rf.APIHandlers.PatchPassword, true, 0))
	r.PATCH("/api/user/password/:username", rf.createAPIHandler(rf.APIHandlers.PatchUserPassword, true, 0))
	r.DELETE("/api/user/:username", rf.createAPIHandler(rf.APIHandlers.DeleteUser, true, 0))

	minClientRank := config.GetPermissionConfig().MinClientRank

	//client routes
	r.GET("/api/clients", rf.createAPIHandler(rf.APIHandlers.GetClients, true, minClientRank))
	r.POST("/api/client", rf.createAPIHandler(rf.APIHandlers.PostClient, true, minClientRank))
	r.PUT("/api/client/:id", rf.createAPIHandler(rf.APIHandlers.PutClient, true, minClientRank))
	r.DELETE("/api/client/:id", rf.createAPIHandler(rf.APIHandlers.DeleteClient, true, minClientRank))

	//user-role routes
	r.GET("/api/client/:id/roles", rf.createAPIHandler(rf.APIHandlers.GetUserRoles, true, 0))
	r.POST("/api/client/:id/role", rf.createAPIHandler(rf.APIHandlers.PostUserRole, true, 0))
	r.PUT("/api/client/:id/role/:username", rf.createAPIHandler(rf.APIHandlers.PutUserRole, true, 0))
	r.DELETE("/api/client/:id/role/:username", rf.createAPIHandler(rf.APIHandlers.DeleteUserRole, true, 0))

	//session routes
	r.POST("/api/session", rf.createAPIHandler(rf.APIHandlers.PostSession, false, 0))
	r.DELETE("/api/session", rf.createAPIHandler(rf.APIHandlers.DeleteSession, true, 0))
}

func (rf CoreRouterFactory) setPanicHandler(r *httprouter.Router, rw writers.ResponseWriter) {
	r.PanicHandler = func(w http.ResponseWriter, _ *http.Request, info interface{}) {
		log.Println(info)
		rw.WriteInternalErrorResponse(w)
	}
}

type handlerFunc func(*http.Request, httprouter.Params, *models.Session, parsers.BodyParser, data.DataCRUD) (int, interface{})

type RouterFactoryConfig struct {
	BodyParser     parsers.BodyParser
	ResponseWriter writers.ResponseWriter
}

func (rf CoreRouterFactory) createUIHandler(handler handlerFunc, authenticateUser bool, minRank int) httprouter.Handle {
	return rf.createHandler(handler, authenticateUser, minRank, RouterFactoryConfig{
		BodyParser:     parsers.NewFormBodyParser(),
		ResponseWriter: writers.NewUIResponseWriter(),
	})
}

func (rf CoreRouterFactory) createAPIHandler(handler handlerFunc, authenticateUser bool, minRank int) httprouter.Handle {
	return rf.createHandler(handler, authenticateUser, minRank, RouterFactoryConfig{
		BodyParser:     parsers.NewJSONBodyParser(),
		ResponseWriter: writers.NewJSONResponseWriter(),
	})
}

func (rf CoreRouterFactory) createHandler(handler handlerFunc, authenticateUser bool, minRank int, cfg RouterFactoryConfig) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var session *models.Session
		var cerr common.CustomError

		err := rf.ScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
			//authenticate the user if required
			if authenticateUser {
				session, cerr = rf.getSession(exec, req)
				if cerr.Type == common.ErrorTypeClient {
					cfg.ResponseWriter.WriteErrorResponse(w, http.StatusUnauthorized, cerr.Error())
					return nil
				}
				if cerr.Type == common.ErrorTypeInternal {
					cfg.ResponseWriter.WriteInternalErrorResponse(w)
					return nil
				}

				//verify the user has at least the min rank required to access the route
				if session.Rank < minRank {
					cfg.ResponseWriter.WriteInsufficientPermissionsErrorResponse(w)
					return nil
				}
			}

			//handle route in transaction scope
			return rf.ScopeFactory.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
				status, res := handler(req, params, session, cfg.BodyParser, tx)

				//handle special redirect case
				if status == http.StatusSeeOther {
					rf.setRedirectHeader(w, res.(common.RedirectResponse))
					cfg.ResponseWriter.WriteResponse(w, status, nil)
					return true, nil
				}

				//send response
				cfg.ResponseWriter.WriteResponse(w, status, res)
				return status == http.StatusOK, nil
			})
		})

		if err != nil {
			log.Println(err)
			cfg.ResponseWriter.WriteInternalErrorResponse(w)
		}
	}
}

func (CoreRouterFactory) setRedirectHeader(w http.ResponseWriter, res common.RedirectResponse) {
	w.Header().Set("Location", res.Location)
	if res.Cookie != "" {
		w.Header().Set("Set-Cookie", res.Cookie)
	}
}

func (CoreRouterFactory) getSession(CRUD models.SessionCRUD, req *http.Request) (*models.Session, common.CustomError) {
	//extract the token string from the authorization header
	splitTokens := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(splitTokens) != 2 {
		return nil, common.ClientError("no bearer token provided")
	}

	//parse the session token
	token, err := uuid.Parse(splitTokens[1])
	if err != nil {
		log.Println(common.ChainError("error parsing token", err))
		return nil, common.ClientError("bearer token was in an invalid format")
	}

	//fetch the session
	session, err := CRUD.GetSessionByToken(token)
	if err != nil {
		log.Println(common.ChainError("error getting session by token", err))
		return nil, common.InternalError()
	}

	//no session found
	if session == nil {
		return nil, common.ClientError("bearer token invalid or expired")
	}

	return session, common.NoError()
}
