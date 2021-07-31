package router_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/router"
	handlermocks "authserver/router/handlers/mocks"
	"authserver/router/mocks"
	testhelpers "authserver/testing"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
	testhelpers.ScopeFactorySuite

	AuthenticatorMock mocks.Authenticator
	HandlersMock      handlermocks.IHandlers
	Router            *httprouter.Router

	Method  string
	Route   string
	Handler string
}

func (suite *RouterTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.AuthenticatorMock = mocks.Authenticator{}
	suite.HandlersMock = handlermocks.IHandlers{}

	suite.AuthenticatorMock.On("Authenticate", mock.Anything, mock.Anything).Return(nil, requesterror.NoError())

	rf := router.RouterFactory{
		Authenticator: &suite.AuthenticatorMock,
		ScopeFactory:  &suite.ScopeFactoryMock,
		Handlers:      &suite.HandlersMock,
	}
	suite.Router = rf.CreateRouter()
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromDataExecutorScope_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromTransactionScope_ReturnsErrorToDataExecutorScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)
	message := "TransactionScope error"

	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(nil, func(err error) {
		//assert
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})
	suite.SetupScopeFactoryMock_CreateTransactionScope(errors.New(message))

	//act
	_, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
}

func (suite *RouterTestSuite) TestRoute_WithNonOKStatusFromHandler_SendsResponseAndReturnsFailureToTransactionScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
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
	common.AssertErrorResponse(&suite.Suite, res, status, message)
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WithOKStatusFromHandler_SendsResponseAndReturnsSuccessToTransactionScope() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
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
	common.AssertSuccessResponse(&suite.Suite, res)
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WhereHandlerPanics_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, suite.Method, server.URL+suite.Route, "", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Run(func(_ mock.Arguments) {
		panic("")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

type RouterAuthTestSuite struct {
	RouterTestSuite
}

func TestPostUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "POST",
		Route:   "/user",
		Handler: "PostUser",
	})
}

func TestDeleteUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "DELETE",
		Route:   "/user",
		Handler: "DeleteUser",
	})
}

func TestPatchUserPasswordTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "PATCH",
		Route:   "/user/password",
		Handler: "PatchUserPassword",
	})
}

func TestPostTokenTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "POST",
		Route:   "/token",
		Handler: "PostToken",
	})
}

func TestDeleteTokenTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "DELETE",
		Route:   "/token",
		Handler: "DeleteToken",
	})
}
