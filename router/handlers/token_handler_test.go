package handlers_test

import (
	"authserver/common"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithClientErrorCreatingTokenRedirectURL_ReturnsBadRequest() {
	//arrange
	body := handlers.PostTokenBody{
		ClientId: uuid.New(),
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create token error"
	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithInternalErrorCreatingTokenRedirectURL_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostTokenBody{
		ClientId: uuid.New(),
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithNoErrors_ReturnsRedirect() {
	//arrange
	body := handlers.PostTokenBody{
		ClientId: uuid.New(),
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	redirectUrl := "redirect.com"
	suite.ControllersMock.On("CreateTokenRedirectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redirectUrl, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostToken(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Equal(http.StatusSeeOther, status)
	suite.Equal(redirectUrl, res)
}

func TestTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TokenHandlerTestSuite{})
}
