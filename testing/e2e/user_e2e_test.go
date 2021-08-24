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
	postUserBody := handlers.PostUserBody{
		Username: username,
		Password: password,
	}
	res := suite.SendRequest(http.MethodPost, "/user", "", postUserBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//login
	token := suite.Login(username, password)

	//update user password
	patchPasswordBody := handlers.PatchUserPasswordBody{
		OldPassword: password,
		NewPassword: "NewPassword123!",
	}
	res = suite.SendRequest(http.MethodPatch, "/user/password", token, patchPasswordBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//delete user
	res = suite.SendRequest(http.MethodDelete, "/user", token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
