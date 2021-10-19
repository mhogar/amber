package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

	suite.RendererMock.AssertCalled(suite.T(), "RenderView", mock.Anything, mock.Anything, "home/index", "partials/page")
}

func TestHomeHandlerTestSuite(t *testing.T) {
	suite.Run(t, &HomeHandlerTestSuite{})
}
