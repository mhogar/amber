package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
)

func (suite *E2ETestSuite) Login(creds UserCredentials) string {
	body := handlers.PostSessionBody{
		Username: creds.Username,
		Password: creds.Password,
	}
	res := suite.SendRequest(http.MethodPost, "/session", "", body)

	return helpers.ParseDataResponseOK(&suite.Suite, res)["token"].(string)
}

func (suite *E2ETestSuite) Logout(token string) {
	res := suite.SendRequest(http.MethodDelete, "/session", token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}
