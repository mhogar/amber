package e2e_test

import (
	"authserver/common"
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

	//set db key and create database
	viper.Set("db_key", "integration")

	//create the test server
	runner := server.CreateHTTPTestServerRunner(dependencies.ResolveRouterFactory())
	suite.Server = runner.Server.(*server.HTTPTestServer).Server

	// run the server
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
	body := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: username,
			Password: password,
			ClientID: config.GetAppId().String(),
		},
	}
	res := suite.SendRequest(http.MethodPost, "/token", "", body)

	tokenRes := common.AccessTokenResponse{}
	helpers.ParseResponseOK(&suite.Suite, res, &tokenRes)

	return tokenRes.AccessToken
}

func (suite *E2ETestSuite) LoginAsMaxAdmin() string {
	return suite.Login("admin", "Admin123!")
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendRequest(http.MethodDelete, "/token", token, nil)
	helpers.ParseAndAssertSuccessResponse(&suite.Suite, res)
}
