package e2e_test

import (
	"net/http"
	"testing"

	"github.com/mhogar/amber/router/handlers"

	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendGetUsersRequest(token string) *http.Response {
	return suite.SendJSONRequest(http.MethodGet, "/users", token, nil)
}

func (suite *E2ETestSuite) SendCreateUserRequest(token string, username string, password string, rank int) *http.Response {
	postUserBody := handlers.PostUserBody{
		Username: username,
		Password: password,
		Rank:     rank,
	}
	return suite.SendJSONRequest(http.MethodPost, "/user", token, postUserBody)
}

func (suite *E2ETestSuite) CreateUser(token string, username string, rank int) UserCredentials {
	newUser := UserCredentials{
		Username: username,
		Password: "Password123!",
	}

	res := suite.SendCreateUserRequest(token, username, newUser.Password, rank)
	suite.ParseAndAssertOKSuccessResponse(res)

	return newUser
}

func (suite *E2ETestSuite) SendUpdateUserRequest(token string, username string, rank int) *http.Response {
	putUserBody := handlers.PutUserBody{
		Rank: rank,
	}
	return suite.SendJSONRequest(http.MethodPut, "/user/"+username, token, putUserBody)
}

func (suite *E2ETestSuite) SendUpdatePasswordRequest(token string, oldPassword string, newPassword string) *http.Response {
	patchPasswordBody := handlers.PatchPasswordBody{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	return suite.SendJSONRequest(http.MethodPatch, "/user/password", token, patchPasswordBody)
}

func (suite *E2ETestSuite) SendUpdateUserPasswordRequest(token string, username string, password string) *http.Response {
	patchUserPasswordBody := handlers.PatchUserPasswordBody{
		Password: password,
	}
	return suite.SendJSONRequest(http.MethodPatch, "/user/password/"+username, token, patchUserPasswordBody)
}

func (suite *E2ETestSuite) SendDeleteUserRequest(token string, username string) *http.Response {
	return suite.SendJSONRequest(http.MethodDelete, "/user/"+username, token, nil)
}

func (suite *E2ETestSuite) DeleteUser(token string, username string) {
	res := suite.SendDeleteUserRequest(token, username)
	suite.ParseAndAssertOKSuccessResponse(res)
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

func (suite *UserE2ETestSuite) TestGetUsers_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendGetUsersRequest("")
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestCreateUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendCreateUserRequest("", "new_user", "Password123!", 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestCreateUser_WithRankLessThanUser_ReturnsForbidden() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "", "", 11)
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)
}

func (suite *UserE2ETestSuite) TestCreateUser_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "", "", 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "username", "cannot be empty")
}

func (suite *UserE2ETestSuite) TestCreateUser_WithNonUniqueUsername_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, suite.ExistingUser.Username, "Password123!", 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "username", "already in use")
}

func (suite *UserE2ETestSuite) TestCreateUser_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	res := suite.SendCreateUserRequest(suite.AdminToken, "new_user", "invalid", 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "password", "does not meet", "criteria")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateUserRequest("", suite.ExistingUser.Username, 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithRankLessThanNewUserRank_ReturnsForbidden() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, 11)
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, "DNE", 0)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithRankLessThanCurrentUserRank_ReturnsForbidden() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user
	res := suite.SendUpdateUserRequest(token, suite.Admin.Username, 0)
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, -1)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "rank", "invalid")
}

func (suite *UserE2ETestSuite) TestUpdateUser_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateUserRequest(suite.AdminToken, suite.ExistingUser.Username, 0)
	suite.ParseAndAssertOKSuccessResponse(res)
}

func (suite *UserE2ETestSuite) TestUpdatePassword_WhereOldPasswordIsIncorrect_ReturnsBadRequest() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdatePasswordRequest(token, "incorrect", "Password1234!")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "old password", "incorrect")

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdatePassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdatePasswordRequest(token, suite.ExistingUser.Password, "invalid")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "password", "does not meet", "criteria")

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdatePassword_WithValidRequest_ReturnsSuccess() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdatePasswordRequest(token, suite.ExistingUser.Password, suite.ExistingUser.Password)
	suite.ParseAndAssertOKSuccessResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateUserPasswordRequest("", "", "")
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendUpdateUserPasswordRequest(suite.AdminToken, "DNE", "")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.ExistingUser)

	//update user password
	res := suite.SendUpdateUserPasswordRequest(token, suite.Admin.Username, "")
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WherePasswordDoesNotMeetCriteria_ReturnsBadRequest() {
	res := suite.SendUpdateUserPasswordRequest(suite.AdminToken, suite.ExistingUser.Username, "invalid")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "password", "does not meet", "criteria")
}

func (suite *UserE2ETestSuite) TestUpdateUserPassword_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateUserPasswordRequest(suite.AdminToken, suite.ExistingUser.Username, suite.ExistingUser.Password)
	suite.ParseAndAssertOKSuccessResponse(res)
}

func (suite *UserE2ETestSuite) TestDeleteUser_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteUserRequest("", suite.ExistingUser.Username)
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserE2ETestSuite) TestDeleteUser_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendDeleteUserRequest(suite.AdminToken, "DNE")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserE2ETestSuite) TestDeleteUser_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.ExistingUser)

	//delete user
	res := suite.SendDeleteUserRequest(token, suite.Admin.Username)
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func TestUserE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserE2ETestSuite{})
}
