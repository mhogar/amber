package handlers_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/models"
	"authserver/router"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenHandlerTestSuite struct {
	RouterTestSuite
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithErrorCreatingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", nil)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", "invalid")

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", "invalid json body")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithMissingGrantType_ReturnsInvalidRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", "missing grant_type parameter")
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithUnsupportedGrantType_ReturnsUnsupportedGrantType() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "unsupported",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "unsupported_grant_type", "")
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithMissingParameters_ReturnsInvalidRequest() {
	var grantBody router.PostTokenPasswordGrantBody
	var expectedErrorDescription string

	testCase := func() {
		//arrange
		server := httptest.NewServer(suite.Router)
		defer server.Close()

		body := router.PostTokenBody{
			GrantType:                  "password",
			PostTokenPasswordGrantBody: grantBody,
		}
		req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

		suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

		//act
		res, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)

		//assert
		suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
		common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_request", expectedErrorDescription)
	}

	grantBody = router.PostTokenPasswordGrantBody{
		Password: "password",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing username parameter"
	suite.Run("MissingUsername", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		ClientID: "client id",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing password parameter"
	suite.Run("MissingPassword", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		Scope:    "scope",
	}
	expectedErrorDescription = "missing client_id parameter"
	suite.Run("MissingClientID", testCase)

	grantBody = router.PostTokenPasswordGrantBody{
		Username: "username",
		Password: "password",
		ClientID: "client id",
	}
	expectedErrorDescription = "missing scope parameter"
	suite.Run("MissingScope", testCase)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithErrorParsingClient_ReturnsInvalidClient() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: "invalid",
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid_client", "client_id", "invalid format")
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithClientErrorCreatingTokenFromPassword_ReturnsInvalidClient() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	errorName := "error_name"
	message := "create token error"
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, requesterror.OAuthClientError(errorName, message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertOAuthErrorResponse(&suite.Suite, res, http.StatusBadRequest, errorName, message)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithInternalErrorCreatingTokenFromPassword_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, requesterror.OAuthInternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestPostToken_WithErrorCommitingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	clientID := uuid.New()
	token := models.CreateNewAccessToken(nil, nil, nil)

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: clientID.String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(token, requesterror.OAuthNoError())
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithValidRequest_ReturnsAccessToken() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	clientID := uuid.New()
	token := models.CreateNewAccessToken(nil, nil, nil)

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: clientID.String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(token, requesterror.OAuthNoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertNotCalled(suite.T(), "Authenticate", mock.Anything)
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "CreateTokenFromPassword", &suite.TransactionMock, body.Username, body.Password, clientID, body.Scope)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	common.AssertAccessTokenResponse(&suite.Suite, res, token.ID.String())
}

func (suite *TokenHandlerTestSuite) TestPostToken_PasswordGrant_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: router.PostTokenPasswordGrantBody{
			Username: "username",
			Password: "password",
			ClientID: uuid.New().String(),
			Scope:    "scope",
		},
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/token", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateTokenFromPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithErrorCreatingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithClientErrorDeletingToken_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	message := "delete token error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithInternalErrorDeletingToken_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithErrorCommitingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteToken", &suite.TransactionMock, token)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *TokenHandlerTestSuite) TestDeleteToken_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/token", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteToken", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func TestTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, &TokenHandlerTestSuite{})
}
