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
	Username string
	ClientID uuid.UUID
}

func (suite *UserRoleE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()

	suite.Username = suite.CreateUser("username", "Password123!", 0)
	suite.ClientID = suite.CreateClient(0, "key.pem")
}

func (suite *UserRoleE2ETestSuite) TearDownSuite() {
	suite.DeleteClient(suite.ClientID)
	suite.DeleteUser(suite.Username)

	suite.E2ETestSuite.TearDownSuite()
}

func (suite *UserRoleE2ETestSuite) TestCreateUserRole_UpdateUserRole_DeleteUserRole() {
	//create user-role
	suite.CreateUserRole(suite.Username, suite.ClientID, "role")

	//update user-role
	putUserRoleBody := handlers.PostUserRoleBody{
		Role: "new role",
	}
	res := suite.SendRequest(http.MethodPut, path.Join("/user", suite.Username, "role", suite.ClientID.String()), suite.AdminToken, putUserRoleBody)
	helpers.ParseDataResponseOK(&suite.Suite, res)

	//delete user-role
	suite.DeleteUserRole(suite.Username, suite.ClientID)
}

func TestUserRoleE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserRoleE2ETestSuite{})
}
