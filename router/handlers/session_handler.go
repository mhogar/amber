package handlers

import (
	"authserver/common"
	"authserver/data"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SessionDataResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// PostSessionBody is the struct the body of requests to PostSession should be parsed into.
type PostSessionBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h CoreHandlers) PostSession(req *http.Request, _ httprouter.Params, _ *models.Session, tx data.Transaction) (int, interface{}) {
	var body PostSessionBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostSession request body", err))
		return common.NewBadRequestResponse("invalid json body")
	}

	//create the session
	session, cerr := h.Controllers.CreateSession(tx, body.Username, body.Password)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return newSessionDataResponse(session)
}

func (h CoreHandlers) DeleteSession(_ *http.Request, _ httprouter.Params, session *models.Session, tx data.Transaction) (int, interface{}) {
	//delete the session
	cerr := h.Controllers.DeleteSession(tx, session)
	if cerr.Type == common.ErrorTypeClient {
		return common.NewBadRequestResponse(cerr.Error())
	}
	if cerr.Type == common.ErrorTypeInternal {
		return common.NewInternalServerErrorResponse()
	}

	return common.NewSuccessResponse()
}

func newSessionDataResponse(session *models.Session) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(SessionDataResponse{
		Token:    session.ID.String(),
		Username: session.User.Username,
	})
}
