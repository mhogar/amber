package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PermissionsE2ETestSuite struct {
	E2ETestSuite
	Username1 string
	Username2 string
	Password  string
	ClientID  uuid.UUID
}

func (suite *PermissionsE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.Password = "Password123!"

	//-- create the users --
	suite.Username1 = "username1"
	suite.CreateUser(suite.Username1, suite.Password, 0)

	suite.Username2 = "username2"
	suite.CreateUser(suite.Username2, suite.Password, 1)

	//create the client
	suite.ClientID = suite.CreateClient(0, "key.pem")
}

func (suite *PermissionsE2ETestSuite) TearDownSuite() {
	suite.DeleteClient(suite.ClientID)

	suite.DeleteUser(suite.Username2)
	suite.DeleteUser(suite.Username1)

	suite.E2ETestSuite.TearDownSuite()
}

func (suite *PermissionsE2ETestSuite) TestCreateUser_WithRankLessThanNewUser_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to create user
	postUserBody := handlers.PostUserBody{
		Username: "username3",
		Password: suite.Password,
		Rank:     1,
	}
	res := suite.SendRequest(http.MethodPost, "/user", token, postUserBody)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestUpdateUser_WithRankLessThanCurrentUserRank_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to update user
	putUserBody := handlers.PutUserBody{
		Rank: 2,
	}
	res := suite.SendRequest(http.MethodPut, "/user/"+suite.Username2, token, putUserBody)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestUpdateUser_WithRankLessThanNewUserRank_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username2, suite.Password)

	//attempt to update user
	putUserBody := handlers.PutUserBody{
		Rank: 2,
	}
	res := suite.SendRequest(http.MethodPut, "/user/"+suite.Username1, token, putUserBody)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestDeleteUser_WithRankLessThanUser_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to delete user
	res := suite.SendRequest(http.MethodDelete, "/user/"+suite.Username2, token, nil)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestCreateClient_WithRankLessThanMin_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to create client
	postClientBody := handlers.PostClientBody{
		Name:        "Test Client",
		RedirectUrl: "https://mhogar.dev",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	res := suite.SendRequest(http.MethodPost, "/client", token, postClientBody)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestUpdateClient_WithRankLessThanMin_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to create client
	putClientBody := handlers.PostClientBody{
		Name:        "Test Client",
		RedirectUrl: "https://mhogar.dev",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	res := suite.SendRequest(http.MethodPut, "/client/"+suite.ClientID.String(), token, putClientBody)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *PermissionsE2ETestSuite) TestDeleteClient_WithRankLessThanMin_ReturnsInsufficientPermissions() {
	//login
	token := suite.Login(suite.Username1, suite.Password)

	//attempt to delete client
	res := suite.SendRequest(http.MethodDelete, "/client/"+suite.ClientID.String(), token, nil)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func TestPermissionsE2ETestSuite(t *testing.T) {
	suite.Run(t, &PermissionsE2ETestSuite{})
}
