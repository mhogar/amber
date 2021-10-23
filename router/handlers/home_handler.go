package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/parsers"
)

func (h CoreUIHandlers) GetHome(req *http.Request, _ httprouter.Params, _ *models.Session, _ parsers.BodyParser, _ data.DataCRUD) (int, interface{}) {
	return h.renderHomeView(req)
}

func (h CoreUIHandlers) renderHomeView(req *http.Request) (int, interface{}) {
	return http.StatusOK, h.Renderer.RenderView(req, nil, "home/index", "partials/page")
}
