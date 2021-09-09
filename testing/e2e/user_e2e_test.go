package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateUserRequest(token string, username string, password string, rank int) *http.Response {
	postUserBody := handlers.PostUserBody{
		Username: username,
		Password: password,
		Rank:     rank,
	}
	return suite.SendRequest(http.MethodPost, "/user", token, postUserBody)
}

func (suite *E2ETestSuite) CreateUser(token string, username string, rank int) UserCredentials {
	newUser := UserCredentials{
		Username: username,
		Password: "Password123!",
	}

	res := suite.SendCreateUserRequest(token, username, newUser.Password, rank)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	return newUser
}

func (suite *E2ETestSuite) SendUpdateUserRequest(token string, username string, rank int) *http.Response {
	putUserBody := handlers.PutUserBody{
		Rank: rank,
	}
	return suite.SendRequest(http.MethodPut, "/user/"+username, token, putUserBody)
}

func (suite *E2ETestSuite) SendUpdateUserPasswordRequest(token string, password string, newPassword string) *http.Response {
	patchPasswordBody := handlers.PatchUserPasswordBody{
		OldPassword: password,
		NewPassword: newPassword,
	}
	return suite.SendRequest(http.MethodPatch, "/user/password", token, patchPasswordBody)
}

func (suite *E2ETestSuite) SendDeleteUserRequest(token string, username string) *http.Response {
	return suite.SendRequest(http.MethodDelete, "/user/"+username, token, nil)
}

func (suite *E2ETestSuite) DeleteUser(token string, username string) {
	res := suite.SendDeleteUserRequest(token, username)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

type UserE2ETestSuite struct {
	E2ETestSuite
	ExistingUser UserCredentials
}

func (suite *UserE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.ExistingUser = suite.CreateUser(suite.AdminToken, "user", 5)
}

func (suite *UserE2ETestSuite) TearDownSuite() {
	suite.DeleteUser(suite.AdminToken, suite.ExistingUser.Username)
	suite.E2ETestSuite.TearDownSuite()
}

func (suite *UserE2ETestSuite) TestCreateUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendCreateUserRequest("", "new_user", "Password123!", 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestCreateUser_WithRankLessThanUser_ReturnsForbidden() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "", "", 11)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserE2ETestSuite) TestCreateUser_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "", "", 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "username", "cannot be empty")
}

func (suite *UserE2ETestSuite) TestCreateUser_WithNonUniqueUsername_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, suite.ExistingUser.Username, "Password123!", 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "username", "already in use")
}

func (suite *UserE2ETestSuite) TestCreateUser_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "new_user", "invalid", 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "password", "does not meet", "criteria")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateUserRequest("", suite.ExistingUser.Username, 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithRankLessThanNewUserRank_ReturnsForbidden() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, 11)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, "DNE", 0)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithRankLessThanCurrentUserRank_ReturnsForbidden() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user
	res := suite.SendUpdateUserRequest(token, suite.Admin.Username, 0)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, -1)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "rank", "invalid")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, 0)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WhereOldPasswordIsIncorrect_ReturnsBadRequest() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdateUserPasswordRequest(token, "incorrect", "Password1234!")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "old password", "incorrect")

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdateUserPasswordRequest(token, suite.ExistingUser.Password, "invalid")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "password", "does not meet", "criteria")

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestDeleteUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteUserRequest("", suite.ExistingUser.Username)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestDeleteUser_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendDeleteUserRequest(suite.AdminToken, "DNE")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserE2ETestSuite) TestDeleteUser_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.ExistingUser)

	//delete user
	res := suite.SendDeleteUserRequest(token, suite.Admin.Username)
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
