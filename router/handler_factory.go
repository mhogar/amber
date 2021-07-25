package router

import (
	requesterror "authserver/common/request_error"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handlerFunc func(*http.Request, httprouter.Params, *models.AccessToken, data.Transaction) (int, interface{})

func (h RouterFactory) createHandler(handler handlerFunc, authenticateUser bool) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var token *models.AccessToken
		var rerr requesterror.RequestError

		err := h.ScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
			//authenticate the user if required
			if authenticateUser {
				token, rerr = h.Authenticator.Authenticate(exec, req)
				if rerr.Type == requesterror.ErrorTypeClient {
					sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
					return nil
				} else if rerr.Type == requesterror.ErrorTypeInternal {
					sendInternalErrorResponse(w)
					return nil
				}
			}

			//handle route in transaction scope
			return h.ScopeFactory.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
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
