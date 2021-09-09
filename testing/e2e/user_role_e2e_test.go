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

type UserRoleE2ETestSuite struct {
	E2ETestSuite
	User1    UserCredentials
	ClientID uuid.UUID
}

func (suite *UserRoleE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.User1 = suite.CreateUser(suite.AdminToken, "user1", 5)
	suite.ClientID = suite.CreateClient(suite.AdminToken, 0, "key.pem")
}

func (suite *UserRoleE2ETestSuite) TearDownSuite() {
	suite.DeleteClient(suite.AdminToken, suite.ClientID)
	suite.DeleteUser(suite.AdminToken, suite.User1.Username)
	suite.E2ETestSuite.TearDownSuite()
}
func (suite *UserRoleE2ETestSuite) TestCreateUserRole_UpdateUserRole_DeleteUserRole() {
	//create user-role
	suite.CreateUserRole(suite.User1.Username, suite.ClientID, "role")

	//update user-role
	putUserRoleBody := handlers.PostUserRoleBody{
		Role: "new role",
	}
	res := suite.SendRequest(http.MethodPut, path.Join("/user", suite.User1.Username, "role", suite.ClientID.String()), suite.AdminToken, putUserRoleBody)
	helpers.ParseDataResponseOK(&suite.Suite, res)

	//delete user-role
	suite.DeleteUserRole(suite.User1.Username, suite.ClientID)
}

func TestUserRoleE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserRoleE2ETestSuite{})
}
