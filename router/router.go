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

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	ResponseTypeRaw  = iota
	ResponseTypeJSON = iota
)

type RouterFactory interface {
	// CreateRouter creates a new httprouter with the endpoints and panic handler configured.
	CreateRouter() *httprouter.Router
}

type CoreRouterFactory struct {
	ScopeFactory data.ScopeFactory
	Handlers     handlers.Handlers
}

func (rf CoreRouterFactory) CreateRouter() *httprouter.Router {
	r := httprouter.New()
	r.PanicHandler = func(w http.ResponseWriter, _ *http.Request, info interface{}) {
		log.Println(info)
		sendInternalErrorResponse(w)
	}

	//user routes
	r.GET("/users", rf.createHandler(rf.Handlers.GetUsers, ResponseTypeJSON, true, 0))
	r.POST("/user", rf.createHandler(rf.Handlers.PostUser, ResponseTypeJSON, true, 0))
	r.PUT("/user/:username", rf.createHandler(rf.Handlers.PutUser, ResponseTypeJSON, true, 0))
	r.PATCH("/user/password", rf.createHandler(rf.Handlers.PatchPassword, ResponseTypeJSON, true, 0))
	r.PATCH("/user/password/:username", rf.createHandler(rf.Handlers.PatchUserPassword, ResponseTypeJSON, true, 0))
	r.DELETE("/user/:username", rf.createHandler(rf.Handlers.DeleteUser, ResponseTypeJSON, true, 0))

	minClientRank := config.GetPermissionConfig().MinClientRank

	//client routes
	r.GET("/clients", rf.createHandler(rf.Handlers.GetClients, ResponseTypeJSON, true, minClientRank))
	r.POST("/client", rf.createHandler(rf.Handlers.PostClient, ResponseTypeJSON, true, minClientRank))
	r.PUT("/client/:id", rf.createHandler(rf.Handlers.PutClient, ResponseTypeJSON, true, minClientRank))
	r.DELETE("/client/:id", rf.createHandler(rf.Handlers.DeleteClient, ResponseTypeJSON, true, minClientRank))

	//user-role routes
	r.GET("/client/:id/roles", rf.createHandler(rf.Handlers.GetUserRoles, ResponseTypeJSON, true, 0))
	r.POST("/client/:id/role", rf.createHandler(rf.Handlers.PostUserRole, ResponseTypeJSON, true, 0))
	r.PUT("/client/:id/role/:username", rf.createHandler(rf.Handlers.PutUserRole, ResponseTypeJSON, true, 0))
	r.DELETE("/client/:id/role/:username", rf.createHandler(rf.Handlers.DeleteUserRole, ResponseTypeJSON, true, 0))

	//session routes
	r.POST("/session", rf.createHandler(rf.Handlers.PostSession, ResponseTypeJSON, false, 0))
	r.DELETE("/session", rf.createHandler(rf.Handlers.DeleteSession, ResponseTypeJSON, true, 0))

	//token routes
	r.GET("/token", rf.createHandler(rf.Handlers.GetToken, ResponseTypeRaw, false, 0))
	r.POST("/token", rf.createHandler(rf.Handlers.PostToken, ResponseTypeRaw, false, 0))

	return r
}

type handlerFunc func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

func (rf CoreRouterFactory) createHandler(handler handlerFunc, responseType int, authenticateUser bool, minRank int) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var session *models.Session
		var cerr common.CustomError

		err := rf.ScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
			//authenticate the user if required
			if authenticateUser {
				session, cerr = rf.getSession(exec, req)
				if cerr.Type == common.ErrorTypeClient {
					sendErrorResponse(w, http.StatusUnauthorized, cerr.Error())
					return nil
				}
				if cerr.Type == common.ErrorTypeInternal {
					sendInternalErrorResponse(w)
					return nil
				}

				//verify the user has at least the min rank required to access the route
				if session.Rank < minRank {
					sendInsufficientPermissionsErrorResponse(w)
					return nil
				}
			}

			//handle route in transaction scope
			return rf.ScopeFactory.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
				status, data := handler(req, params, session, tx)

				//handle special redirect case
				if status == http.StatusSeeOther {
					w.Header().Set("Location", data.(string))
					sendRawResponse(w, status, nil)
					return true, nil
				}

				//send response based on type (default to raw)
				if responseType == ResponseTypeJSON {
					sendJSONResponse(w, status, data)
				} else {
					sendRawResponse(w, status, data.([]byte))
				}

				return status == http.StatusOK, nil
			})
		})

		if err != nil {
			log.Println(err)
			sendInternalErrorResponse(w)
		}
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
