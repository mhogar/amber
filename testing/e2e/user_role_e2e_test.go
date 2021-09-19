package e2e_test

import (
	"net/http"
	"path"
	"testing"

	"github.com/mhogar/amber/router/handlers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendGetUserRolesRequest(token string, clientID string) *http.Response {
	return suite.SendJSONRequest(http.MethodGet, path.Join("/client", clientID, "roles"), token, nil)
}

func (suite *E2ETestSuite) SendCreateUserRoleRequest(token string, clientID string, username string, role string) *http.Response {
	postUserRoleBody := handlers.PostUserRoleBody{
		Username: username,
		Role:     role,
	}
	return suite.SendJSONRequest(http.MethodPost, path.Join("/client", clientID, "role"), token, postUserRoleBody)
}

func (suite *E2ETestSuite) CreateUserRole(token string, clientID uuid.UUID, username string, role string) {
	res := suite.SendCreateUserRoleRequest(token, clientID.String(), username, role)
	suite.ParseAndAssertOKSuccessResponse(res)
}

func (suite *E2ETestSuite) SendUpdateUserRoleRequest(token string, clientID string, username string, role string) *http.Response {
	putUserRoleBody := handlers.PostUserRoleBody{
		Role: role,
	}
	return suite.SendJSONRequest(http.MethodPut, path.Join("/client", clientID, "role", username), token, putUserRoleBody)
}

func (suite *E2ETestSuite) SendDeleteUserRoleRequest(token string, clientID string, username string) *http.Response {
	return suite.SendJSONRequest(http.MethodDelete, path.Join("/client", clientID, "role", username), token, nil)
}

func (suite *E2ETestSuite) DeleteUserRole(token string, clientID uuid.UUID, username string) {
	res := suite.SendDeleteUserRoleRequest(token, clientID.String(), username)
	suite.ParseAndAssertOKSuccessResponse(res)
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
	suite.CreateUserRole(suite.AdminToken, suite.ClientID, suite.User.Username, "role")
}

func (suite *UserRoleE2ETestSuite) TearDownSuite() {
	suite.DeleteUserRole(suite.AdminToken, suite.ClientID, suite.User.Username)
	suite.DeleteClient(suite.AdminToken, suite.ClientID)
	suite.DeleteUser(suite.AdminToken, suite.User.Username)

	suite.E2ETestSuite.TearDownSuite()
}

func (suite *UserRoleE2ETestSuite) TestGetUserRoles_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendGetUserRolesRequest("", suite.ClientID.String())
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestGetUserRoles_WithInvalidClientId_ReturnsBadRequest() {
	res := suite.SendGetUserRolesRequest(suite.AdminToken, "invalid")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendCreateUserRoleRequest("", suite.ClientID.String(), suite.User.Username, "role")
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithInvalidClientID_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, "invalid", suite.User.Username, "role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), "DNE", "role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendCreateUserRoleRequest(token, suite.ClientID.String(), suite.Admin.Username, "role")
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), suite.User.Username, "")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "role", "cannot be empty")
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_WhereRoleForClientAlreadyExists_ReturnsBadRequest() {
	res := suite.SendCreateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), suite.User.Username, "role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "already has a role", "client")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendUpdateUserRoleRequest("", suite.ClientID.String(), suite.User.Username, "new role")
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidClientID_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, "invalid", suite.User.Username, "new role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), "DNE", "new role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendUpdateUserRoleRequest(token, suite.ClientID.String(), suite.Admin.Username, "new role")
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithInvalidBody_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), suite.User.Username, "")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "role", "cannot be empty")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WhereRoleNotFound_ReturnsBadRequest() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, uuid.New().String(), suite.User.Username, "new role")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "no role found", "user", "client")
}

func (suite *UserRoleE2ETestSuite) TestUpdateUserRole_WithValidRequest_ReturnsSuccess() {
	res := suite.SendUpdateUserRoleRequest(suite.AdminToken, suite.ClientID.String(), suite.User.Username, "new role")
	suite.ParseAndAssertOKSuccessResponse(res)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithInvalidSession_ReturnsUnauthorized() {
	res := suite.SendDeleteUserRoleRequest("", suite.ClientID.String(), suite.User.Username)
	suite.ParseAndAssertErrorResponse(res, http.StatusUnauthorized)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithInvalidClientID_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, "invalid", suite.User.Username)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "client id", "invalid format")
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WhereUsernameDoesNotExist_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, suite.ClientID.String(), "DNE")
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "user", "not found")
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WithRankLessThanUser_ReturnsForbidden() {
	//login
	token := suite.Login(suite.User)

	//delete user
	res := suite.SendDeleteUserRoleRequest(token, suite.ClientID.String(), suite.Admin.Username)
	suite.ParseAndAssertInsufficientPermissionsErrorResponse(res)

	//logout
	suite.Logout(token)
}

func (suite *UserRoleE2ETestSuite) TestDeleteUserRole_WhereRoleNotFound_ReturnsBadRequest() {
	res := suite.SendDeleteUserRoleRequest(suite.AdminToken, uuid.New().String(), suite.User.Username)
	suite.ParseAndAssertErrorResponse(res, http.StatusBadRequest, "no role found", "user", "client")
}

func TestUserRoleE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserRoleE2ETestSuite{})
}
