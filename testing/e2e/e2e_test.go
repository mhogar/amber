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

type E2ETestSuite struct {
	suite.Suite
	Server     *httptest.Server
	AdminToken string
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
	suite.AdminToken = suite.Login("admin", "Admin123!")
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

func (suite *E2ETestSuite) Login(username string, password string) string {
	body := handlers.PostSessionBody{
		Username: username,
		Password: password,
	}
	res := suite.SendRequest(http.MethodPost, "/session", "", body)

	return helpers.ParseDataResponseOK(&suite.Suite, res)["token"].(string)
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendRequest(http.MethodDelete, "/session", token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) CreateUser(username string, password string, rank int) string {
	postUserBody := handlers.PostUserBody{
		Username: username,
		Password: password,
		Rank:     rank,
	}
	res := suite.SendRequest(http.MethodPost, "/user", suite.AdminToken, postUserBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	return username
}

func (suite *E2ETestSuite) DeleteUser(username string) {
	res := suite.SendRequest(http.MethodDelete, "/user/"+username, suite.AdminToken, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) CreateClient(tokenType int, keyUri string) uuid.UUID {
	postClientBody := handlers.PostClientBody{
		Name:        "Test Client",
		RedirectUrl: "https://mhogar.dev",
		TokenType:   tokenType,
		KeyUri:      keyUri,
	}
	res := suite.SendRequest(http.MethodPost, "/client", suite.AdminToken, postClientBody)

	id, err := uuid.Parse(helpers.ParseDataResponseOK(&suite.Suite, res)["id"].(string))
	suite.Require().NoError(err)

	return id
}

func (suite *E2ETestSuite) DeleteClient(id uuid.UUID) {
	res := suite.SendRequest(http.MethodDelete, "/client/"+id.String(), suite.AdminToken, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
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
