package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserRoleCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRolesForClient_WithInvalidUserRoles_ReturnsError() {
	//arrange
	roles := make([]*models.UserRole, 1)
	roles[0] = models.CreateUserRole("", uuid.Nil, "")

	//act
	err := suite.Tx.UpdateUserRolesForClient(uuid.New(), roles)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRolesForClient_UpdatesRolesForClient() {
	//arrange
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 0, []byte("password")))
	user3 := suite.SaveUser(models.CreateUser("user3", 0, []byte("password")))
	user4 := suite.SaveUser(models.CreateUser("user4", 0, []byte("password")))

	//-- first update --
	roles := make([]*models.UserRole, 2)
	roles[0] = models.CreateUserRole(user1.Username, client.UID, "role1")
	roles[1] = models.CreateUserRole(user2.Username, client.UID, "role2")

	//act
	err := suite.Tx.UpdateUserRolesForClient(client.UID, roles)
	suite.Require().NoError(err)

	//assert
	res, err := suite.Tx.GetUserRolesForClient(client.UID)
	suite.Require().NoError(err)
	suite.Equal(roles, res)

	//-- second update --
	roles[0] = models.CreateUserRole(user3.Username, client.UID, "role1")
	roles[1] = models.CreateUserRole(user4.Username, client.UID, "role2")

	//act
	err = suite.Tx.UpdateUserRolesForClient(client.UID, roles)
	suite.Require().NoError(err)

	//assert
	res, err = suite.Tx.GetUserRolesForClient(client.UID)
	suite.Require().NoError(err)
	suite.Equal(roles, res)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleForClient_TestCases() {
	//arrange
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 0, []byte("password")))

	roles := make([]*models.UserRole, 1)
	roles[0] = models.CreateUserRole(user2.Username, client.UID, "role")

	err := suite.Tx.UpdateUserRolesForClient(client.UID, roles)
	suite.Require().NoError(err)

	clientUID := uuid.New()
	username := ""
	var expectedRole *models.UserRole = nil

	testCase := func() {
		//act
		role, err := suite.Tx.GetUserRoleForClient(clientUID, username)

		//assert
		suite.Require().NoError(err)
		suite.Equal(expectedRole, role)
	}

	//-- test cases --
	suite.Run("ClientNotFound_ReturnsNilRole", testCase)

	clientUID = client.UID
	suite.Run("UserNotFound_ReturnsNilRole", testCase)

	username = user1.Username
	suite.Run("NoRoleForUserAndClient_ReturnsNilRole", testCase)

	username = user2.Username
	expectedRole = roles[0]
	suite.Run("WithRoleForUserAndClient_ReturnsRole", testCase)
}

func TestUserRoleCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleCRUDTestSuite{})
}
