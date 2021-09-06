package e2e_test

import (
	"authserver/models"
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
	Password string
	ClientID uuid.UUID
}

func (suite *UserRoleE2ETestSuite) SetupTest() {
	suite.Username = "username"
	suite.Password = "Password123!"

	//create new user and login
	suite.CreateUser(suite.Username, suite.Password)
	suite.Login(suite.Username, suite.Password)

	//create new client
	suite.ClientID = suite.CreateClient(models.ClientTokenTypeDefault, "keys/test.private.pem")
}

func (suite *UserRoleE2ETestSuite) TearDownTest() {
	suite.DeleteClient(suite.ClientID)
	suite.DeleteUser()
}

func (suite *UserRoleE2ETestSuite) Test_UpdateUserRolesForClient() {
	//-- update user roles --
	rolesBody := make([]handlers.PutClientRolesBody, 1)
	rolesBody[0] = handlers.PutClientRolesBody{
		Username: suite.Username,
		Role:     "role",
	}

	res := suite.SendRequest(http.MethodPut, path.Join("/client", suite.ClientID.String(), "roles"), suite.Token, rolesBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func TestUserRoleE2ETestSuite(t *testing.T) {
	suite.Run(t, &UserRoleE2ETestSuite{})
}
