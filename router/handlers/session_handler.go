package handlers

import (
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"

	"github.com/julienschmidt/httprouter"
)

type SessionDataResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type PostSessionBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h CoreHandlers) PostSession(req *http.Request, _ httprouter.Params, _ *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	var body PostSessionBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostSession request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the session
	session, cerr := h.Controllers.CreateSession(CRUD, body.Username, body.Password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return h.newSessionDataResponse(session)
}

func (h CoreHandlers) DeleteSession(_ *http.Request, _ httprouter.Params, session *models.Session, CRUD data.DataCRUD) (int, interface{}) {
	//delete the session
	cerr := h.Controllers.DeleteSession(CRUD, session.Token)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func (CoreHandlers) newSessionDataResponse(session *models.Session) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(SessionDataResponse{
		Token:    session.Token.String(),
		Username: session.Username,
	})
}
