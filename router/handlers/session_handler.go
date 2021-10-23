package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/parsers"

	"github.com/julienschmidt/httprouter"
)

type LoginViewData struct {
	Error string
}

type SessionDataResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type PostSessionBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h CoreUIHandlers) GetLogin(req *http.Request, _ httprouter.Params, session *models.Session, _ parsers.BodyParser, _ data.DataCRUD) (int, interface{}) {
	//redirect to home if already logged-in
	if session != nil {
		return common.NewRedirectResponse(getBaseURL(req, "home"), "")
	}

	return h.renderLoginView(req, "")
}

func (h CoreUIHandlers) PostSession(req *http.Request, params httprouter.Params, session *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	//proxy to API handler
	status, res := h.API.PostSession(req, params, session, parser, CRUD)
	if status != http.StatusOK {
		return h.renderLoginView(req, res.(common.ErrorResponse).Error)
	}

	//redirect to home on success
	token := res.(common.DataResponse).Data.(SessionDataResponse).Token
	return common.NewRedirectResponse(getBaseURL(req, "home"), fmt.Sprintf("token=%s; SameSite=Strict; HttpOnly", token))
}

func (h CoreUIHandlers) renderLoginView(req *http.Request, err string) (int, interface{}) {
	data := LoginViewData{
		Error: err,
	}

	return http.StatusOK, h.Renderer.RenderView(req, data, "login/index", "partials/login_form", "partials/alert")
}

func (h CoreAPIHandlers) PostSession(req *http.Request, _ httprouter.Params, _ *models.Session, parser parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
	var body PostSessionBody

	//parse the body
	err := parser.ParseBody(req, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostSession request body", err))
		return common.NewBadRequestResponse("invalid request body")
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

func (h CoreAPIHandlers) DeleteSession(_ *http.Request, _ httprouter.Params, session *models.Session, _ parsers.BodyParser, CRUD data.DataCRUD) (int, interface{}) {
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

func (CoreAPIHandlers) newSessionDataResponse(session *models.Session) (int, common.DataResponse) {
	return common.NewSuccessDataResponse(SessionDataResponse{
		Token:    session.Token.String(),
		Username: session.Username,
	})
}
