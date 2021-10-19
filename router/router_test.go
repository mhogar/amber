package router_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router"
	handlermocks "github.com/mhogar/amber/router/handlers/mocks"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const MinClientRank = 5

type RouterTestSuite struct {
	helpers.ScopeFactorySuite
	HandlersMock handlermocks.Handlers
	Router       *httprouter.Router
	Server       *httptest.Server

	Method       string
	Route        string
	Handler      string
	ResponseType int

	Session *models.Session
	TokenId string
}

func (suite *RouterTestSuite) SetupSuite() {
	viper.Set("permission", config.PermissionConfig{
		MinClientRank: MinClientRank,
	})
}

func (suite *RouterTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.HandlersMock = handlermocks.Handlers{}

	suite.Session = nil
	suite.TokenId = ""

	rf := router.CoreRouterFactory{
		ScopeFactory: &suite.ScopeFactoryMock,
		Handlers:     &suite.HandlersMock,
	}
	suite.Router = rf.CreateRouter()
	suite.Server = httptest.NewServer(suite.Router)
}

func (suite *RouterTestSuite) TearDownTest() {
	suite.Server.Close()
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromDataExecutorScope_ReturnsInternalServerError() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertInternalServerErrorResponse(res)
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromTransactionScope_ReturnsErrorToDataExecutorScope() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)
	message := "TransactionScope error"

	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(nil, func(err error) {
		//assert
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(errors.New(message))

	//act
	_, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
}

func (suite *RouterTestSuite) TestRoute_WithNonOKStatusFromHandler_SendsResponseAndReturnsFailureToTransactionScope() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.False(result)
		suite.NoError(err)
	})

	status := http.StatusInternalServerError
	message := "error response"

	var body interface{}
	if suite.ResponseType == router.ResponseTypeJSON {
		body = common.NewErrorResponse(message)
	} else {
		body = []byte(message)
	}

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(status, body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	if suite.ResponseType == router.ResponseTypeJSON {
		suite.ParseAndAssertErrorResponse(res, status, message)
	} else {
		suite.ReadAndAssertRawResponse(res, status, body.([]byte))
	}
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WithRedirectStatusFromHandler_SendsRedirectResponseAndReturnsSuccessToTransactionScope() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.True(result)
		suite.NoError(err)
	})

	status := http.StatusSeeOther
	redirectUrl := "https://mhogar.dev"
	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(status, redirectUrl)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.Equal(redirectUrl, res.Request.URL.String())
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WithOKStatusFromHandler_SendsResponseAndReturnsSuccessToTransactionScope() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.True(result)
		suite.NoError(err)
	})

	status := http.StatusOK
	var body interface{}
	if suite.ResponseType == router.ResponseTypeJSON {
		_, body = common.NewSuccessResponse()
	} else {
		body = []byte("handler result")
	}

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(status, body)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	if suite.ResponseType == router.ResponseTypeJSON {
		suite.ParseAndAssertOKSuccessResponse(res)
	} else {
		suite.ReadAndAssertRawResponse(res, status, body.([]byte))
	}
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WhereHandlerPanics_ReturnsInternalServerError() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	suite.HandlersMock.On(suite.Handler, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Run(func(_ mock.Arguments) {
		panic("")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertInternalServerErrorResponse(res)
}

type RouterAuthTestSuite struct {
	RouterTestSuite
	MinRank int
}

func (suite *RouterAuthTestSuite) SetupTest() {
	suite.RouterTestSuite.SetupTest()

	token := uuid.New()
	suite.Session = models.CreateSession(token, "username", suite.MinRank)
	suite.TokenId = token.String()
}

func (suite *RouterAuthTestSuite) TestRoute_WithNoBearerToken_ReturnsUnauthorized() {
	//arrange
	var req *http.Request

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	testCase := func() {
		//act
		res, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)

		//assert
		suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized, "no bearer token")
	}

	req = suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, "", nil)
	suite.Run("NoAuthorizationHeader", testCase)

	req.Header.Set("Authorization", "invalid")
	suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
}

func (suite *RouterAuthTestSuite) TestRoute_WithBearerTokenInInvalidFormat_ReturnsUnauthorized() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, "invalid", nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized, "bearer token", "invalid format")
}

func (suite *RouterAuthTestSuite) TestRoute_WithErrorGettingSessionByID_ReturnsInternalServerError() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(nil, errors.New(""))
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertInternalServerErrorResponse(res)
}

func (suite *RouterAuthTestSuite) TestRoute_WhereSessionWithIDisNotFound_ReturnsUnauthorized() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(nil, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized, "bearer token", "invalid", "expired")
}

func (suite *RouterAuthTestSuite) TestRoute_WithSessionRankLessThanMinRank_ReturnsForbidden() {
	//arrange
	req := suite.CreateJSONRequest(suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)
	session := models.CreateNewSession("username", suite.MinRank-1)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)
}

func TestGetHomeTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:       "GET",
		Route:        "/",
		Handler:      "GetHome",
		ResponseType: router.ResponseTypeRaw,
	})
}

func TestGetUsersTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "GET",
			Route:        "/users",
			Handler:      "GetUsers",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPostUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "POST",
			Route:        "/user",
			Handler:      "PostUser",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPutUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "PUT",
			Route:        "/user/username",
			Handler:      "PutUser",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPatchPasswordTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "PATCH",
			Route:        "/user/password",
			Handler:      "PatchPassword",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPatchUserPasswordTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "PATCH",
			Route:        "/user/password/username",
			Handler:      "PatchUserPassword",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestDeleteUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "DELETE",
			Route:        "/user/username",
			Handler:      "DeleteUser",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestGetClientsTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "GET",
			Route:        "/clients",
			Handler:      "GetClients",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: MinClientRank,
	})
}

func TestPostClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "POST",
			Route:        "/client",
			Handler:      "PostClient",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: MinClientRank,
	})
}

func TestPutClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "PUT",
			Route:        "/client/0",
			Handler:      "PutClient",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: MinClientRank,
	})
}

func TestDeleteClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "DELETE",
			Route:        "/client/0",
			Handler:      "DeleteClient",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: MinClientRank,
	})
}

func TestGetUserRolesTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "GET",
			Route:        "/client/0/roles",
			Handler:      "GetUserRoles",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPostUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "POST",
			Route:        "/client/0/role",
			Handler:      "PostUserRole",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPutUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "PUT",
			Route:        "/client/0/role/username",
			Handler:      "PutUserRole",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestDeleteUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "DELETE",
			Route:        "/client/0/role/username",
			Handler:      "DeleteUserRole",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestPostSessionTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:       "POST",
		Route:        "/session",
		Handler:      "PostSession",
		ResponseType: router.ResponseTypeJSON,
	})
}

func TestDeleteSessionTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:       "DELETE",
			Route:        "/session",
			Handler:      "DeleteSession",
			ResponseType: router.ResponseTypeJSON,
		},
		MinRank: 0,
	})
}

func TestGetTokenTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:       "GET",
		Route:        "/token",
		Handler:      "GetToken",
		ResponseType: router.ResponseTypeRaw,
	})
}

func TestPostTokenTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:       "POST",
		Route:        "/token",
		Handler:      "PostToken",
		ResponseType: router.ResponseTypeRaw,
	})
}
