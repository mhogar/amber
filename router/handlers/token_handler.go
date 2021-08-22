package handlers

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PostTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h CoreHandlers) PostToken(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PostTokenBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostToken request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the token
	cerr := h.Controllers.CreateToken(CRUD, body.Username, body.Password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}
