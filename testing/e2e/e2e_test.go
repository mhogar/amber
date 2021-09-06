package e2e_test

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/models"
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
	Server *httptest.Server
	Token  string
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

func (suite *E2ETestSuite) Login(username string, password string) {
	body := handlers.PostSessionBody{
		Username: username,
		Password: password,
	}
	res := suite.SendRequest(http.MethodPost, "/session", "", body)

	suite.Token = helpers.ParseDataResponseOK(&suite.Suite, res)["token"].(string)
}

func (suite *E2ETestSuite) LoginAsMaxAdmin() {
	suite.Login("admin", "Admin123!")
}

func (suite *E2ETestSuite) Logout() {
	res := suite.SendRequest(http.MethodDelete, "/session", suite.Token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) CreateUser(username string, password string) {
	postUserBody := handlers.PostUserBody{
		Username: username,
		Password: password,
	}
	res := suite.SendRequest(http.MethodPost, "/user", "", postUserBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) DeleteUser() {
	res := suite.SendRequest(http.MethodDelete, "/user", suite.Token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) CreateClient(tokenType int, keyUri string) uuid.UUID {
	postClientBody := handlers.PostClientBody{
		Name:        "Test Client",
		RedirectUrl: "https://mhogar.dev",
		TokenType:   tokenType,
		KeyUri:      keyUri,
	}
	res := suite.SendRequest(http.MethodPost, "/client", suite.Token, postClientBody)

	id, err := uuid.Parse(helpers.ParseDataResponseOK(&suite.Suite, res)["id"].(string))
	suite.Require().NoError(err)

	return id
}

func (suite *E2ETestSuite) DeleteClient(id uuid.UUID) {
	res := suite.SendRequest(http.MethodDelete, "/client/"+id.String(), suite.Token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) UpdateUserRolesForClient(clientID uuid.UUID, roles ...*models.UserRole) []*models.UserRole {
	rolesBody := make([]handlers.PutClientRolesBody, len(roles))
	for index, role := range roles {
		rolesBody[index] = handlers.PutClientRolesBody{
			Username: role.Username,
			Role:     role.Role,
		}
	}

	res := suite.SendRequest(http.MethodPut, path.Join("/client", clientID.String(), "roles"), suite.Token, rolesBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	return roles
}
