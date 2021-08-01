package handlers_test

import (
	"authserver/common"
	"authserver/models"
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

func (suite *TokenHandlerTestSuite) TestPostToken_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertOAuthErrorResponse(&suite.Suite, res, "invalid_request", "invalid json body")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithMissingGrantType_ReturnsInvalidRequest() {
	//arrange
	body := handlers.PostTokenBody{}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertOAuthErrorResponse(&suite.Suite, res, "invalid_request", "missing grant_type parameter")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithUnsupportedGrantType_ReturnsUnsupportedGrantType() {
	//arrange
	body := handlers.PostTokenBody{
		GrantType: "unsupported",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertOAuthErrorResponse(&suite.Suite, res, "unsupported_grant_type", "")
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithMissingParameters_ReturnsInvalidRequest() {
	var grantBody handlers.PostTokenPasswordGrantBody
	var expectedErrorDescription string

	testCase := func() {
		//arrange
		body := handlers.PostTokenBody{
			GrantType:                  "password",
			PostTokenPasswordGrantBody: grantBody,
		}
		req := helpers.CreateDummyRequest(&suite.Suite, body)

		//act
		status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

		//assert
		suite.Equal(http.StatusBadRequest, status)
		helpers.AssertOAuthErrorResponse(&suite.Suite, res, "invalid_request", expectedErrorDescription)
	}

	grantBody = handlers.PostTokenPasswordGrantBody{
		Password: "password",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing username parameter"
	suite.Run("MissingUsername", testCase)

	grantBody = handlers.PostTokenPasswordGrantBody{
		Username: "username",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing password parameter"
	suite.Run("MissingPassword", testCase)

	grantBody = handlers.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing client_id parameter"
	suite.Run("MissingClientID", testCase)

	grantBody = handlers.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		ClientID: "client id",
	}
	expectedErrorDescription = "missing scope parameter"
	suite.Run("MissingScope", testCase)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithErrorParsingClient_ReturnsInvalidClient() {
	//arrange
	body := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: "invalid",
			Scope:    "scope",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertOAuthErrorResponse(&suite.Suite, res, "invalid_client", "client_id", "invalid format")
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithClientErrorCreatingTokenFromPassword_ReturnsInvalidClient() {
	//arrange
	body := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	errorName := "error_name"
	message := "create token error"
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.OAuthClientError(errorName, message))

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertOAuthErrorResponse(&suite.Suite, res, errorName, message)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithInternalErrorCreatingTokenFromPassword_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.OAuthInternalError())

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithValidRequest_ReturnsAccessToken() {
	//arrange
	clientID := uuid.New()
	token := models.CreateNewAccessToken(nil, nil, nil)

	body := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: clientID.String(),
			Scope:    "scope",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(token, common.OAuthNoError())

	//act
	status, res := suite.Handlers.PostToken(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusOK, status)
	helpers.AssertAccessTokenResponse(&suite.Suite, res, token.ID.String())

	suite.ControllersMock.AssertCalled(suite.T(), "CreateTokenFromPassword", &suite.TransactionMock, body.Username, body.Password, clientID, body.Scope)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithClientErrorDeletingToken_ReturnsBadRequest() {
	//arrange
	token := &models.AccessToken{}

	message := "delete token error"
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.Handlers.DeleteToken(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithInternalErrorDeletingToken_ReturnsInternalServerError() {
	//arrange
	token := &models.AccessToken{}

	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.Handlers.DeleteToken(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithValidRequest_ReturnsSuccess() {
	//arrange
	token := &models.AccessToken{}

	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.Handlers.DeleteToken(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteToken", &suite.TransactionMock, token)
}

func TestTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TokenHandlerTestSuite{})
}