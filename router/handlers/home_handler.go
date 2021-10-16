package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
)

func (h CoreHandlers) GetHome(req *http.Request, _ httprouter.Params, _ *models.Session, _ data.DataCRUD) (int, interface{}) {
	return h.renderHomeView(req)
}

func (h CoreHandlers) renderHomeView(req *http.Request) (int, interface{}) {
	return http.StatusOK, h.Renderer.RenderView(req, nil, "home/index", "partials/navbar")
}
