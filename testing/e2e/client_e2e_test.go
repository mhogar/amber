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

func (suite *ClientE2ETestSuite) TestCreateClient_UpdateClient_DeleteClient() {
	//create client
	id := suite.CreateClient(0, "key.pem")

	//update client
	putClientBody := handlers.PostClientBody{
		Name:        "New Name",
		RedirectUrl: "redirect2.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	res := suite.SendRequest(http.MethodPut, "/client/"+id.String(), suite.AdminToken, putClientBody)
	helpers.ParseDataResponseOK(&suite.Suite, res)

	//delete client
	suite.DeleteClient(id)
}

func TestClientE2ETestSuite(t *testing.T) {
	suite.Run(t, &ClientE2ETestSuite{})
}
