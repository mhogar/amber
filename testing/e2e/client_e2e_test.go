package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendGetClientsRequest(token string) *http.Response {
	return suite.SendJSONRequest(http.MethodGet, "/clients", token, nil)
}

func (suite *E2ETestSuite) SendCreateClientRequest(token string, tokenType int, keyUri string) *http.Response {
	postClientBody := handlers.PostClientBody{
		Name:        "Test Client",
		RedirectUrl: "https://mhogar.dev",
		TokenType:   tokenType,
		KeyUri:      keyUri,
	}
	return suite.SendJSONRequest(http.MethodPost, "/client", token, postClientBody)
}

func (suite *E2ETestSuite) CreateClient(token string, tokenType int, keyUri string) uuid.UUID {
	res := suite.SendCreateClientRequest(token, tokenType, keyUri)

	id, err := uuid.Parse(helpers.ParseDataResponseOK(&suite.Suite, res)["id"].(string))
	suite.Require().NoError(err)

	return id
}

func (suite *E2ETestSuite) SendUpdateClientRequest(token string, id string, tokenType int, keyUri string) *http.Response {
	putClientBody := handlers.PostClientBody{
		Name:        "New Name",
		RedirectUrl: "redirect2.com",
		TokenType:   tokenType,
		KeyUri:      keyUri,
	}
	return suite.SendJSONRequest(http.MethodPut, "/client/"+id, token, putClientBody)
}

func (suite *E2ETestSuite) SendDeleteClientRequest(token string, id string) *http.Response {
	return suite.SendJSONRequest(http.MethodDelete, "/client/"+id, token, nil)
}

func (suite *E2ETestSuite) DeleteClient(token string, id uuid.UUID) {
	res := suite.SendDeleteClientRequest(token, id.String())
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

type ClientE2ETestSuite struct {
	E2ETestSuite
	User     UserCredentials
	ClientId uuid.UUID
}

func (suite *ClientE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.User = suite.CreateUser(suite.AdminToken, "user", 0)
	suite.ClientId = suite.CreateClient(suite.AdminToken, 0, "key.pem")
}

func (suite *ClientE2ETestSuite) TearDownSuite() {
	suite.DeleteClient(suite.AdminToken, suite.ClientId)
	suite.DeleteUser(suite.AdminToken, suite.User.Username)
	suite.E2ETestSuite.TearDownSuite()
}

func (suite *ClientE2ETestSuite) TestGetClients_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendGetClientsRequest("")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *ClientE2ETestSuite) TestGetClients_WithRankLessThanMin_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//get clients
	res := suite.SendGetClientsRequest(token)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *ClientE2ETestSuite) TestCreateClient_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendCreateClientRequest("", 0, "key.pem")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *ClientE2ETestSuite) TestCreateClient_WithRankLessThanMin_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//create client
	res := suite.SendCreateClientRequest(token, 0, "key.pem")
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *ClientE2ETestSuite) TestCreateClient_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendCreateClientRequest(suite.AdminToken, -1, "key.pem")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "token type", "invalid")
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateClientRequest("", suite.ClientId.String(), 0, "key.pem")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WithRankLessThanMin_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//create client
	res := suite.SendUpdateClientRequest(token, suite.ClientId.String(), 0, "key.pem")
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WithInvalidClientId_ReturnsBadRequest() {
	res := suite.SendUpdateClientRequest(suite.AdminToken, "invalid", 0, "key.pem")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendUpdateClientRequest(suite.AdminToken, suite.ClientId.String(), -1, "")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "token type", "invalid")
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WhereClientNotFound_ReturnsBadRequest() {
	res := suite.SendUpdateClientRequest(suite.AdminToken, uuid.New().String(), 0, "key.pem")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client", "not found")
}

func (suite *ClientE2ETestSuite) TestUpdateClient_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateClientRequest(suite.AdminToken, suite.ClientId.String(), 0, "key.pem")
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *ClientE2ETestSuite) TestDeleteClient_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteClientRequest("", suite.ClientId.String())
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *ClientE2ETestSuite) TestDeleteClient_WithRankLessThanMin_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete client
	res := suite.SendDeleteClientRequest(token, suite.ClientId.String())
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *ClientE2ETestSuite) TestDeleteClient_WithInvalidClientId_ReturnsBadRequest() {
	res := suite.SendDeleteClientRequest(suite.AdminToken, "invalid")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *ClientE2ETestSuite) TestDeleteClient_WhereClientNotFound_ReturnsBadRequest() {
	res := suite.SendDeleteClientRequest(suite.AdminToken, uuid.New().String())
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client", "not found")
}

func TestClientE2ETestSuite(t *testing.T) {
	suite.Run(t, &ClientE2ETestSuite{})
}
