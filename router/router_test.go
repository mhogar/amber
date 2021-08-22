package router_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router"
	handlermocks "authserver/router/handlers/mocks"
	"authserver/testing/helpers"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	helpers.ScopeFactorySuite
	HandlersMock handlermocks.Handlers
	Router       *httprouter.Router

	Method  string
	Route   string
	Handler string
	TokenId string
}

func (suite *RouterTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.HandlersMock = handlermocks.Handlers{}

	suite.TokenId = ""

	rf := router.CoreRouterFactory{
		CoreScopeFactory: &suite.ScopeFactoryMock,
		Handlers:         &suite.HandlersMock,
	}
	suite.Router = rf.CreateRouter()
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromDataExecutorScope_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromTransactionScope_ReturnsErrorToDataExecutorScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)
	message := "TransactionScope error"

	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(nil, func(err error) {
		//assert
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(&models.Session{}, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(errors.New(message))

	//act
	_, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
}

func (suite *RouterTestSuite) TestRoute_WithNonOKStatusFromHandler_SendsResponseAndReturnsFailureToTransactionScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(&models.Session{}, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.False(result)
		suite.NoError(err)
	})

	status := http.StatusInternalServerError
	message := "error response"
	body := common.NewErrorResponse(message)

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(status, body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, status, message)
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WithOKStatusFromHandler_SendsResponseAndReturnsSuccessToTransactionScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(&models.Session{}, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.True(result)
		suite.NoError(err)
	})

	status, body := common.NewSuccessResponse()
	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(status, body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertSuccessResponse(&suite.Suite, res)
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WhereHandlerPanics_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(&models.Session{}, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Run(func(_ mock.Arguments) {
		panic("")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
}

type RouterAuthTestSuite struct {
	RouterTestSuite
}

func (suite *RouterAuthTestSuite) SetupTest() {
	suite.RouterTestSuite.SetupTest()
	suite.TokenId = uuid.New().String()
}

func (suite *RouterAuthTestSuite) TestRoute_WithNoBearerToken_ReturnsUnauthorized() {
	//arrange
	var req *http.Request

	server := httptest.NewServer(suite.Router)
	defer server.Close()

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	testCase := func() {
		//act
		res, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)

		//assert
		helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, "no bearer token")
	}

	req = helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)
	suite.Run("NoAuthorizationHeader", testCase)

	req.Header.Set("Authorization", "invalid")
	suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
}

func (suite *RouterAuthTestSuite) TestRoute_WithBearerTokenInInvalidFormat_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "invalid", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, "bearer token", "invalid format")
}

func (suite *RouterAuthTestSuite) TestRoute_WithErrorGettingSessionByID_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(nil, errors.New(""))
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *RouterAuthTestSuite) TestRoute_WhereSessionWithIDisNotFound_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := helpers.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByID", mock.Anything).Return(nil, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, "bearer token", "invalid", "expired")
}

func TestPostUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "POST",
		Route:   "/user",
		Handler: "PostUser",
	})
}

func TestDeleteUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "DELETE",
			Route:   "/user",
			Handler: "DeleteUser",
		},
	})
}

func TestPatchUserPasswordTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "PATCH",
			Route:   "/user/password",
			Handler: "PatchUserPassword",
		},
	})
}

func TestPostClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "POST",
			Route:   "/client",
			Handler: "PostClient",
		},
	})
}

func TestPutClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "PUT",
			Route:   "/client/0",
			Handler: "PutClient",
		},
	})
}

func TestDeleteClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "DELETE",
			Route:   "/client/0",
			Handler: "DeleteClient",
		},
	})
}

func TestPostSessionTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "POST",
		Route:   "/session",
		Handler: "PostSession",
	})
}

func TestDeleteSessionTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite{
			Method:  "DELETE",
			Route:   "/session",
			Handler: "DeleteSession",
		},
	})
}
