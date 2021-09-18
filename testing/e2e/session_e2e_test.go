package e2e_test

import (
	"authserver/router/handlers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateSessionRequest(username string, password string) *http.Response {
	body := handlers.PostSessionBody{
		Username: username,
		Password: password,
	}
	return suite.SendJSONRequest(http.MethodPost, "/session", "", body)
}

func (suite *E2ETestSuite) Login(creds UserCredentials) string {
	res := suite.SendCreateSessionRequest(creds.Username, creds.Password)
	return suite.ParseDataResponseOK(res)["token"].(string)
}

func (suite *E2ETestSuite) SendDeleteSessionRequest(token string) *http.Response {
	return suite.SendJSONRequest(http.MethodDelete, "/session", token, nil)
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendDeleteSessionRequest(token)
	suite.ParseAndAssertOKSuccessResponse(res)
}

type SessionE2ETestSuite struct {
	E2ETestSuite
}

func (suite *SessionE2ETestSuite) TestCreateSession_WhereUserNotFound_ReturnsBadRequest() {
	res := suite.SendCreateSessionRequest("DNE", "")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "invalid username and/or password")
}

func (suite *SessionE2ETestSuite) TestCreateSession_WithIncorrectPassword_ReturnsBadRequest() {
	res := suite.SendCreateSessionRequest(suite.Admin.Username, "incorrect")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "invalid username and/or password")
}

func (suite *SessionE2ETestSuite) TestDeleteSession_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteSessionRequest("")
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func TestSessionE2ETestSuite(t *testing.T) {
	suite.Run(t, &SessionE2ETestSuite{})
}
