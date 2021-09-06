package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserE2ETestSuite struct {
	E2ETestSuite
}

func (suite *UserE2ETestSuite) TestCreateUser_UpdateUser_DeleteUser() {
	username := "username"
	password := "Password123!"

	//create user
	suite.CreateUser(username, password, 0)

	//update user
	putUserBody := handlers.PutUserBody{
		Rank: 1,
	}
	res := suite.SendRequest(http.MethodPut, "/user/"+username, suite.AdminToken, putUserBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//delete user
	suite.DeleteUser(username)
}

func (suite *UserE2ETestSuite) TestCreateUser_Login_UpdatePassword_Logout_DeleteUser() {
	username := "username"
	password := "Password123!"

	//create user
	suite.CreateUser(username, password, 0)

	//login
	token := suite.Login(username, password)

	//update user password
	patchPasswordBody := handlers.PatchUserPasswordBody{
		OldPassword: password,
		NewPassword: "NewPassword123!",
	}
	res := suite.SendRequest(http.MethodPatch, "/user/password", token, patchPasswordBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)

	//delete user
	suite.DeleteUser(username)
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
