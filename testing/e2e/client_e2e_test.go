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
	suite.LoginAsMaxAdmin()

	//create client
	id := suite.CreateClient(0, "key.pem")

	//update client
	putClientBody := handlers.PostClientBody{
		Name:        "New Name",
		RedirectUrl: "redirect2.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	res := suite.SendRequest(http.MethodPut, "/client/"+id.String(), suite.Token, putClientBody)
	helpers.ParseDataResponseOK(&suite.Suite, res)

	//delete client
	suite.DeleteClient(id)

	//logout
	suite.Logout()
}

func TestClientE2ETestSuite(t *testing.T) {
	suite.Run(t, &ClientE2ETestSuite{})
}
