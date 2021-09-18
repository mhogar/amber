package e2e_test

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/server"
	"authserver/testing/helpers"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/spf13/viper"
)

type UserCredentials struct {
	Username string
	Password string
}

type E2ETestSuite struct {
	helpers.CustomSuite
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

func (suite *E2ETestSuite) SendRequest(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	return res
}

func (suite *E2ETestSuite) SendJSONRequest(method string, endpoint string, bearerToken string, body interface{}) *http.Response {
	return suite.SendRequest(suite.CreateJSONRequest(method, suite.Server.URL+endpoint, bearerToken, body))
}

func (suite *E2ETestSuite) SendFormRequest(method string, endpoint string, bearerToken string, body url.Values) *http.Response {
	return suite.SendRequest(suite.CreateFormRequest(method, suite.Server.URL+endpoint, bearerToken, body))
}
