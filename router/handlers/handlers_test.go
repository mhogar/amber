package handlers_test

import (
	controllermocks "authserver/controllers/mocks"
	datamocks "authserver/data/mocks"
	"authserver/router/handlers"
	renderermocks "authserver/router/renderer/mocks"
	"authserver/testing/helpers"

	"github.com/stretchr/testify/mock"
)

type HandlersTestSuite struct {
	helpers.CustomSuite
	CRUDMock        datamocks.DataCRUD
	ControllersMock controllermocks.Controllers
	RendererMock    renderermocks.Renderer
	CoreHandlers    handlers.CoreHandlers

	RenderViewResult []byte
	RenderViewData   interface{}
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
	suite.ControllersMock = controllermocks.Controllers{}
	suite.RendererMock = renderermocks.Renderer{}

	suite.RenderViewResult = []byte("render view result")
	suite.RendererMock.On("RenderView", mock.Anything, mock.Anything).Return(suite.RenderViewResult).Run(func(args mock.Arguments) {
		suite.RenderViewData = args.Get(1)
	})

	suite.CoreHandlers = handlers.CoreHandlers{
		Controllers: &suite.ControllersMock,
		Renderer:    &suite.RendererMock,
	}
}

func (suite *HandlersTestSuite) AssertRenderViewResult(expected interface{}) {
	suite.Equal(expected, suite.RenderViewResult)
}
