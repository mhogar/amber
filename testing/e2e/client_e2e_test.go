package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientE2ETestSuite struct {
	E2ETestSuite
}

func (suite *ClientE2ETestSuite) TestLogin_CreateClient_UpdateClient_DeleteClient_Logout() {
	//login
	token := suite.LoginAsMaxAdmin()

	//create client
	postClientBody := handlers.PostClientBody{
		Name:        "Name",
		RedirectUrl: "redirect.com",
	}
	res := suite.SendRequest(http.MethodPost, "/client", token, postClientBody)
	id := helpers.ParseDataResponseOK(&suite.Suite, res)["id"].(string)

	//update client
	putClientBody := handlers.PutClientBody{
		Name:        "New Name",
		RedirectUrl: "redirect2.com",
	}
	res = suite.SendRequest(http.MethodPut, "/client/"+id, token, putClientBody)
	helpers.ParseDataResponseOK(&suite.Suite, res)

	//delete client
	res = suite.SendRequest(http.MethodDelete, "/client/"+id, token, nil)
	helpers.ParseAndAssertSuccessResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func TestClientE2ETestSuite(t *testing.T) {
	suite.Run(t, &ClientE2ETestSuite{})
}
