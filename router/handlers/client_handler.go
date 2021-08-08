package handlers

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h CoreHandlers) PostClient(req *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	return common.NewSuccessResponse()
}

func (h CoreHandlers) PutClient(req *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	return common.NewSuccessResponse()
}

func (h CoreHandlers) DeleteClient(_ *http.Request, _ httprouter.Params, token *models.AccessToken, tx data.Transaction) (int, interface{}) {
	return common.NewSuccessResponse()
}
