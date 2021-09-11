package handlers_test

import (
	"authserver/common"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *TokenHandlerTestSuite) SetupSuite() {
	viper.Set("app_name", "App Name")
}

func (suite *TokenHandlerTestSuite) AssertTokenViewRenderedWithData(clientID string, errSubStrings ...string) {
	suite.RendererMock.AssertCalled(suite.T(), "RenderView", "token.gohtml", mock.Anything)

	data := suite.RenderViewData.(handlers.TokenViewData)
	suite.Equal(viper.GetString("app_name"), data.AppName)
	suite.Equal(clientID, data.ClientID)
	helpers.AssertContainsSubstrings(&suite.Suite, data.Error, errSubStrings...)
}

func (suite *TokenHandlerTestSuite) TestGetToken_RendersTokenView() {
	//arrange
	clientID := uuid.New().String()
	req := helpers.CreateRequest(&suite.Suite, "", "/token?client_id="+clientID, "", nil)

	//act
	status, res := suite.CoreHandlers.GetToken(req, nil, nil, nil)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)
	suite.AssertTokenViewRenderedWithData(clientID)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithErrorParsingClientId_RendersTokenViewWithError() {
	//arrange
	clientID := "invalid"
	values := url.Values{
		"client_id": []string{clientID},
		"username":  []string{"username"},
		"password":  []string{"password"},
	}
	req := helpers.CreateDummyFormRequest(&suite.Suite, values)

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)
	suite.AssertTokenViewRenderedWithData(clientID, "client_id", "not provided", "invalid format")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithClientErrorCreatingTokenRedirectURL_RendersTokenViewWithError() {
	//arrange
	clientID := uuid.New().String()
	values := url.Values{
		"client_id": []string{clientID},
		"username":  []string{"username"},
		"password":  []string{"password"},
	}
	req := helpers.CreateDummyFormRequest(&suite.Suite, values)

	message := "create token error"
	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)
	suite.AssertTokenViewRenderedWithData(clientID, message)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithInternalErrorCreatingTokenRedirectURL_RendersTokenViewWithError() {
	//arrange
	clientID := uuid.New().String()
	values := url.Values{
		"client_id": []string{clientID},
		"username":  []string{"username"},
		"password":  []string{"password"},
	}
	req := helpers.CreateDummyFormRequest(&suite.Suite, values)

	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.AssertRenderViewResult(res)
	suite.AssertTokenViewRenderedWithData(clientID, "internal error")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithNoErrors_ReturnsRedirect() {
	//arrange
	clientID := uuid.New().String()
	values := url.Values{
		"client_id": []string{clientID},
		"username":  []string{"username"},
		"password":  []string{"password"},
	}
	req := helpers.CreateDummyFormRequest(&suite.Suite, values)

	redirectUrl := "redirect.com"
	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redirectUrl, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusSeeOther, status)
	suite.Equal(redirectUrl, res)
}

func TestTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TokenHandlerTestSuite{})
}
