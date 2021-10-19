package handlers_test

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HomeHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *HomeHandlerTestSuite) TestGetHome_RendersHomeView() {
	//arrange
	req := suite.CreateRequest("", "/", "", nil)

	//act
	status, res := suite.CoreHandlers.GetHome(req, nil, nil, nil)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)

	suite.RendererMock.AssertCalled(suite.T(), "RenderView", mock.Anything, mock.Anything, "home/index")
}
