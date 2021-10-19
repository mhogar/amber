package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
)

type LoginViewData struct {
	Error string
}

func (h CoreHandlers) GetLogin(req *http.Request, _ httprouter.Params, _ *models.Session, _ data.DataCRUD) (int, interface{}) {
	return h.renderLoginView(req)
}

func (h CoreHandlers) renderLoginView(req *http.Request) (int, interface{}) {
	data := LoginViewData{}

	return http.StatusOK, h.Renderer.RenderView(req, data, "login/index", "partials/login_form", "partials/alert")
}
