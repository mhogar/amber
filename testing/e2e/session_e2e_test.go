package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateSessionRequest(username string, password string) *http.Response {
	body := handlers.PostSessionBody{
		Username: username,
		Password: password,
	}
	return suite.SendRequest(http.MethodPost, "/session", "", body)
}

func (suite *E2ETestSuite) Login(creds UserCredentials) string {
	res := suite.SendCreateSessionRequest(creds.Username, creds.Password)
	return helpers.ParseDataResponseOK(&suite.Suite, res)["token"].(string)
}

func (suite *E2ETestSuite) SendDeleteSessionRequest(token string) *http.Response {
	return suite.SendRequest(http.MethodDelete, "/session", token, nil)
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendDeleteSessionRequest(token)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

type SessionE2ETestSuite struct {
	E2ETestSuite
}

func (suite *SessionE2ETestSuite) TestCreateSession_WhereUserNotFound_ReturnsBadRequest() {
	res := suite.SendCreateSessionRequest("DNE", "")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid username and/or password")
}

func (suite *SessionE2ETestSuite) TestCreateSession_WithIncorrectPassword_ReturnsBadRequest() {
	res := suite.SendCreateSessionRequest(suite.Admin.Username, "incorrect")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid username and/or password")
}

func (suite *SessionE2ETestSuite) TestDeleteSession_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteSessionRequest("")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func TestSessionE2ETestSuite(t *testing.T) {
	suite.Run(t, &SessionE2ETestSuite{})
}
