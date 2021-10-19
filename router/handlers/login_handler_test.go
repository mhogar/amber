package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *LoginHandlerTestSuite) TestGetLogin_RendersLoginView() {
	//arrange
	req := suite.CreateRequest("", "/", "", nil)

	//act
	status, res := suite.CoreHandlers.GetLogin(req, nil, nil, nil)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)

	suite.RendererMock.AssertCalled(suite.T(), "RenderView", mock.Anything, mock.Anything, "login/index", "partials/login_form", "partials/alert")
}

func TestLoginHandlerTestSuite(t *testing.T) {
	suite.Run(t, &LoginHandlerTestSuite{})
}
