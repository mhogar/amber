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
	Server       *httptest.Server

	Method  string
	Route   string
	Handler string

	Session *models.Session
	TokenId string
}

func (suite *RouterTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.HandlersMock = handlermocks.Handlers{}

	suite.Session = nil
	suite.TokenId = ""

	rf := router.CoreRouterFactory{
		CoreScopeFactory: &suite.ScopeFactoryMock,
		Handlers:         &suite.HandlersMock,
	}
	suite.Router = rf.CreateRouter()
	suite.Server = httptest.NewServer(suite.Router)
}

func (suite *RouterTestSuite) TearDownTest() {
	suite.Server.Close()
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromDataExecutorScope_ReturnsInternalServerError() {
	//arrange
	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *RouterTestSuite) TestRoute_WithErrorFromTransactionScope_ReturnsErrorToDataExecutorScope() {
	//arrange
	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)
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

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
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

func (suite *RouterTestSuite) TestRoute_WithRedirectStatusFromHandler_SendsRedirectResponseAndReturnsSuccessToTransactionScope() {
	//arrange

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

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

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(suite.Session, nil)
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
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
	suite.HandlersMock.AssertCalled(suite.T(), suite.Handler, mock.Anything, mock.Anything, mock.Anything, &suite.TransactionMock)
}

func (suite *RouterTestSuite) TestRoute_WhereHandlerPanics_ReturnsInternalServerError() {
	//arrange

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

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
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
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
		helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, "no bearer token")
	}

	req = helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, "", nil)
	suite.Run("NoAuthorizationHeader", testCase)

	req.Header.Set("Authorization", "invalid")
	suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
}

func (suite *RouterAuthTestSuite) TestRoute_WithBearerTokenInInvalidFormat_ReturnsUnauthorized() {
	//arrange

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, "invalid", nil)

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

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(nil, errors.New(""))
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *RouterAuthTestSuite) TestRoute_WhereSessionWithIDisNotFound_ReturnsUnauthorized() {
	//arrange

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(nil, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, "bearer token", "invalid", "expired")
}

func (suite *RouterAuthTestSuite) TestRoute_WithSessionRankLessThanMinRank_ReturnsForbidden() {
	//arrange

	req := helpers.CreateRequest(&suite.Suite, suite.Method, suite.Server.URL+suite.Route, suite.TokenId, nil)
	session := models.CreateNewSession("username", suite.MinRank-1)

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.DataExecutorMock.On("GetSessionByToken", mock.Anything).Return(session, nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func TestPostUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "POST",
			Route:   "/user",
			Handler: "PostUser",
		},
		MinRank: 0,
	})
}

func TestPutUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "PUT",
			Route:   "/user/username",
			Handler: "PutUser",
		},
		MinRank: 0,
	})
}

func TestPatchUserPasswordTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "PATCH",
			Route:   "/user/password",
			Handler: "PatchUserPassword",
		},
		MinRank: 0,
	})
}

func TestDeleteUserTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "DELETE",
			Route:   "/user/username",
			Handler: "DeleteUser",
		},
		MinRank: 0,
	})
}

func TestPostClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "POST",
			Route:   "/client",
			Handler: "PostClient",
		},
		MinRank: 1,
	})
}

func TestPutClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "PUT",
			Route:   "/client/0",
			Handler: "PutClient",
		},
		MinRank: 1,
	})
}

func TestDeleteClientTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "DELETE",
			Route:   "/client/0",
			Handler: "DeleteClient",
		},
		MinRank: 1,
	})
}

func TestPostUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "POST",
			Route:   "/user/username/role",
			Handler: "PostUserRole",
		},
		MinRank: 0,
	})
}

func TestPutUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "PUT",
			Route:   "/user/username/role/0",
			Handler: "PutUserRole",
		},
		MinRank: 0,
	})
}

func TestDeleteUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &RouterAuthTestSuite{
		RouterTestSuite: RouterTestSuite{
			Method:  "DELETE",
			Route:   "/user/username/role/0",
			Handler: "DeleteUserRole",
		},
		MinRank: 0,
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
		RouterTestSuite: RouterTestSuite{
			Method:  "DELETE",
			Route:   "/session",
			Handler: "DeleteSession",
		},
		MinRank: 0,
	})
}

func TestPostTokenTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{
		Method:  "POST",
		Route:   "/token",
		Handler: "PostToken",
	})
}
