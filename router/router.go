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
	// CreateRouter creates a new httprouter with the endpoints and panic handler configured
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
	r.POST("/user", rf.createHandler(rf.Handlers.PostUser, false))
	r.DELETE("/user", rf.createHandler(rf.Handlers.DeleteUser, true))
	r.PATCH("/user/password", rf.createHandler(rf.Handlers.PatchUserPassword, true))

	//client routes
	r.POST("/client", rf.createHandler(rf.Handlers.PostClient, true))
	r.PUT("/client", rf.createHandler(rf.Handlers.PutClient, true))
	r.DELETE("/client", rf.createHandler(rf.Handlers.DeleteClient, true))

	//token routes
	r.POST("/token", rf.createHandler(rf.Handlers.PostToken, false))
	r.DELETE("/token", rf.createHandler(rf.Handlers.DeleteToken, true))

	return r
}

type handlerFunc func(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

func (rf CoreRouterFactory) createHandler(handler handlerFunc, authenticateUser bool) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var token *models.AccessToken
		var rerr common.CustomError

		err := rf.CoreScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
			//authenticate the user if required
			if authenticateUser {
				token, rerr = rf.getAccessToken(exec, req)
				if rerr.Type == common.ErrorTypeClient {
					sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
					return nil
				} else if rerr.Type == common.ErrorTypeInternal {
					sendInternalErrorResponse(w)
					return nil
				}
			}

			//handle route in transaction scope
			return rf.CoreScopeFactory.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
				status, body := handler(req, params, token, tx)
				sendResponse(w, status, body)

				return status == http.StatusOK, nil
			})
		})

		if err != nil {
			log.Println(err)
			sendInternalErrorResponse(w)
		}
	}
}

func (rf CoreRouterFactory) getAccessToken(CRUD models.AccessTokenCRUD, req *http.Request) (*models.AccessToken, common.CustomError) {
	//extract the token string from the authorization header
	splitTokens := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(splitTokens) != 2 {
		return nil, common.ClientError("no bearer token provided")
	}

	//parse the token
	tokenID, err := uuid.Parse(splitTokens[1])
	if err != nil {
		log.Println(common.ChainError("error parsing access token id", err))
		return nil, common.ClientError("bearer token was in an invalid format")
	}

	//fetch the token
	token, err := CRUD.GetAccessTokenByID(tokenID)
	if err != nil {
		log.Println(common.ChainError("error getting access token by id", err))
		return nil, common.InternalError()
	}

	// no token found
	if token == nil {
		return nil, common.ClientError("bearer token invalid or expired")
	}

	// auth success
	return token, common.NoError()
}
