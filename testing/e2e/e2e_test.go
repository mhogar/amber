package e2e_test

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/router/handlers"
	"authserver/server"
	"authserver/testing/helpers"
	"net/http"
	"net/http/httptest"
	"path"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type UserCredentials struct {
	Username string
	Password string
}

type E2ETestSuite struct {
	suite.Suite
	Server *httptest.Server

	AdminToken string
	Admin      UserCredentials
}

func (suite *E2ETestSuite) SetupSuite() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	//set db key
	viper.Set("db_key", "integration")

	//create the test server
	runner := server.CreateHTTPTestServerRunner(dependencies.ResolveRouterFactory())
	suite.Server = runner.Server.(*server.HTTPTestServer).Server

	//run the server
	err = runner.Run()
	suite.Require().NoError(err)

	//login as the max admin
	suite.Admin = UserCredentials{
		Username: "admin",
		Password: "Admin123!",
	}
	suite.AdminToken = suite.Login(suite.Admin)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.Logout(suite.AdminToken)
	suite.Server.Close()
}

func (suite *E2ETestSuite) SendRequest(method string, endpoint string, bearerToken string, body interface{}) *http.Response {
	req := helpers.CreateRequest(&suite.Suite, method, suite.Server.URL+endpoint, bearerToken, body)

	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	return res
}

func (suite *E2ETestSuite) CreateUserRole(username string, clientID uuid.UUID, role string) {
	postUserRoleBody := handlers.PostUserRoleBody{
		ClientID: clientID,
		Role:     role,
	}
	res := suite.SendRequest(http.MethodPost, path.Join("/user", username, "role"), suite.AdminToken, postUserRoleBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) DeleteUserRole(username string, clientID uuid.UUID) {
	res := suite.SendRequest(http.MethodDelete, path.Join("/user", username, "role", clientID.String()), suite.AdminToken, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}
