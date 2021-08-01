package e2e_test

import (
	"authserver/common"
	"authserver/config"
	"authserver/router/handlers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserE2ETestSuite struct {
	E2ETestSuite
}

func (suite *UserE2ETestSuite) TearDownSuite() {
	//close server
	suite.Server.Close()
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
	common.ParseAndAssertSuccessResponse(&suite.Suite, res)

	//login
	postTokenBody := handlers.PostTokenBody{
		GrantType: "password",
		PostTokenPasswordGrantBody: handlers.PostTokenPasswordGrantBody{
			Username: username,
			Password: password,
			ClientID: config.GetAppId().String(),
			Scope:    "all",
		},
	}
	res = suite.SendRequest(http.MethodPost, "/token", "", postTokenBody)

	tokenRes := common.AccessTokenResponse{}
	common.ParseAndAssertResponseOK(&suite.Suite, res, &tokenRes)

	//update user password
	patchBody := handlers.PatchUserPasswordBody{
		OldPassword: password,
		NewPassword: "NewPassword123!",
	}
	res = suite.SendRequest(http.MethodPatch, "/user/password", tokenRes.AccessToken, patchBody)
	common.ParseAndAssertSuccessResponse(&suite.Suite, res)

	//delete user
	res = suite.SendRequest(http.MethodDelete, "/user", tokenRes.AccessToken, nil)
	common.ParseAndAssertSuccessResponse(&suite.Suite, res)
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
