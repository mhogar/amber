package router

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"authserver/router/handlers"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type RouterFactory interface {
	// CreateRouter creates a new httprouter with the endpoints and panic handler configured.
	CreateRouter() *httprouter.Router
}

type CoreRouterFactory struct {
	CoreScopeFactory data.ScopeFactory
	Handlers         handlers.Handlers
}

func (rf CoreRouterFactory) CreateRouter() *httprouter.Router {
	r := httprouter.New()
	r.PanicHandler = panicHandler

	//user routes
	r.POST("/user", rf.createHandler(rf.Handlers.PostUser, true))
	r.PUT("/user/:username", rf.createHandler(rf.Handlers.PutUser, true))
	r.PATCH("/user/password", rf.createHandler(rf.Handlers.PatchUserPassword, true))
	r.DELETE("/user/:username", rf.createHandler(rf.Handlers.DeleteUser, true))

	//client routes
	r.POST("/client", rf.createHandler(rf.Handlers.PostClient, true))
	r.PUT("/client/:id", rf.createHandler(rf.Handlers.PutClient, true))
	r.DELETE("/client/:id", rf.createHandler(rf.Handlers.DeleteClient, true))

	//user-role routes
	r.PUT("/client/:id/roles", rf.createHandler(rf.Handlers.PutClientRoles, true))

	//session routes
	r.POST("/session", rf.createHandler(rf.Handlers.PostSession, false))
	r.DELETE("/session", rf.createHandler(rf.Handlers.DeleteSession, true))

	//token routes
	r.POST("/token", rf.createHandler(rf.Handlers.PostToken, false))

	return r
}

type handlerFunc func(*http.Request, httprouter.Params, *models.Session, data.DataCRUD) (int, interface{})

func (rf CoreRouterFactory) createHandler(handler handlerFunc, authenticateUser bool) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var session *models.Session
		var cerr common.CustomError

		err := rf.CoreScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
			//authenticate the user if required
			if authenticateUser {
				session, cerr = rf.getSession(exec, req)
				if cerr.Type == common.ErrorTypeClient {
					sendErrorResponse(w, http.StatusUnauthorized, cerr.Error())
					return nil
				} else if cerr.Type == common.ErrorTypeInternal {
					sendInternalErrorResponse(w)
					return nil
				}
			}

			//handle route in transaction scope
			return rf.CoreScopeFactory.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
				status, data := handler(req, params, session, tx)

				//handle special redirect case
				if status == http.StatusSeeOther {
					w.Header().Set("Location", data.(string))
					sendResponse(w, status, nil)
					return true, nil
				}

				sendResponse(w, status, data)
				return status == http.StatusOK, nil
			})
		})

		if err != nil {
			log.Println(err)
			sendInternalErrorResponse(w)
		}
	}
}

func (rf CoreRouterFactory) getSession(CRUD models.SessionCRUD, req *http.Request) (*models.Session, common.CustomError) {
	//extract the token string from the authorization header
	splitTokens := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(splitTokens) != 2 {
		return nil, common.ClientError("no bearer token provided")
	}

	//parse the session token
	token, err := uuid.Parse(splitTokens[1])
	if err != nil {
		log.Println(common.ChainError("error parsing token", err))
		return nil, common.ClientError("bearer session was in an invalid format")
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
