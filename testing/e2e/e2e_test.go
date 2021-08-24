package e2e_test

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/router/handlers"
	"authserver/server"
	"authserver/testing/helpers"
	"net/http"
	"net/http/httptest"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	Server *httptest.Server
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
}

func (suite *E2ETestSuite) TearDownSuite() {
	//close server
	suite.Server.Close()
}

func (suite *E2ETestSuite) SendRequest(method string, endpoint string, bearerToken string, body interface{}) *http.Response {
	req := helpers.CreateRequest(&suite.Suite, method, suite.Server.URL+endpoint, bearerToken, body)

	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	return res
}

func (suite *E2ETestSuite) Login(username string, password string) string {
	body := handlers.PostSessionBody{
		Username: username,
		Password: password,
	}
	res := suite.SendRequest(http.MethodPost, "/session", "", body)

	return helpers.ParseDataResponseOK(&suite.Suite, res)["token"].(string)
}

func (suite *E2ETestSuite) LoginAsMaxAdmin() string {
	return suite.Login("admin", "Admin123!")
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendRequest(http.MethodDelete, "/session", token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}
