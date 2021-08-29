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

func (suite *UserE2ETestSuite) TestCreateUser_Login_UpdateUserPassword_DeleteUser() {
	username := "username"
	password := "Password123!"

	//create user
	suite.CreateUser(username, password)

	//login
	suite.Login(username, password)

	//update user password
	patchPasswordBody := handlers.PatchUserPasswordBody{
		OldPassword: password,
		NewPassword: "NewPassword123!",
	}
	res := suite.SendRequest(http.MethodPatch, "/user/password", suite.Token, patchPasswordBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//delete user
	suite.DeleteUser()
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
