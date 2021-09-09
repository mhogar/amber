package e2e_test

import (
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateUserRoleRequest(token string, username string, clientID uuid.UUID, role string) *http.Response {
	postUserRoleBody := handlers.PostUserRoleBody{
		ClientID: clientID,
		Role:     role,
	}
	return suite.SendRequest(http.MethodPost, path.Join("/user", username, "role"), token, postUserRoleBody)
}

func (suite *E2ETestSuite) CreateUserRole(token string, username string, clientID uuid.UUID, role string) {
	res := suite.SendCreateUserRoleRequest(token, username, clientID, role)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *E2ETestSuite) SendUpdateUserRoleRequest(token string, username string, clientID string, role string) *http.Response {
	putUserRoleBody := handlers.PostUserRoleBody{
		Role: role,
	}
	return suite.SendRequest(http.MethodPut, path.Join("/user", username, "role", clientID), token, putUserRoleBody)
}

func (suite *E2ETestSuite) SendDeleteUserRoleRequest(token string, username string, clientID string) *http.Response {
	return suite.SendRequest(http.MethodDelete, path.Join("/user", username, "role", clientID), token, nil)
}

func (suite *E2ETestSuite) DeleteUserRole(token string, username string, clientID uuid.UUID) {
	res := suite.SendDeleteUserRoleRequest(token, username, clientID.String())
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

type UserRoleE2ETestSuite struct {
	E2ETestSuite
	User     UserCredentials
	ClientID uuid.UUID
}

func (suite *UserRoleE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()

	suite.User = suite.CreateUser(suite.AdminToken, "user", 5)
	suite.ClientID = suite.CreateClient(suite.AdminToken, 0, "key.pem")
	suite.CreateUserRole(suite.AdminToken, suite.User.Username, suite.ClientID, "role")
}

func (suite *UserRoleE2ETestSuite) TearDownSuite() {
	suite.DeleteUserRole(suite.AdminToken, suite.User.Username, suite.ClientID)
	suite.DeleteClient(suite.AdminToken, suite.ClientID)
	suite.DeleteUser(suite.AdminToken, suite.User.Username)

	suite.E2ETestSuite.TearDownSuite()
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendCreateUserRoleRequest("", suite.User.Username, suite.ClientID, "role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, "DNE", suite.ClientID, "role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendCreateUserRoleRequest(token, suite.Admin.Username, suite.ClientID, "role")
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, suite.User.Username, suite.ClientID, "")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "role", "cannot be empty")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WhereRoleForClientAlreadyExists_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, suite.User.Username, suite.ClientID, "role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "already has a role", "client")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateUserRoleRequest("", suite.User.Username, suite.ClientID.String(), "new role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidClientID_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.User.Username, "invalid", "new role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, "DNE", suite.ClientID.String(), "new role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendUpdateUserRoleRequest(token, suite.Admin.Username, suite.ClientID.String(), "new role")
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.User.Username, suite.ClientID.String(), "")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "role", "cannot be empty")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WhereRoleNotFound_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.User.Username, uuid.New().String(), "new role")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "no role found", "user", "client")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.User.Username, suite.ClientID.String(), "new role")
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteUserRoleRequest("", suite.User.Username, suite.ClientID.String())
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithInvalidClientID_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, suite.User.Username, "invalid")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, "DNE", suite.ClientID.String())
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendDeleteUserRoleRequest(token, suite.Admin.Username, suite.ClientID.String())
	helpers.ParseAndAssertInsufficientPermissionsErrorResponse(&suite.Suite, res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WhereRoleNotFound_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, suite.User.Username, uuid.New().String())
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "no role found", "user", "client")
}

func TestUserRoleE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserRoleE2ETestSuite{})
}
