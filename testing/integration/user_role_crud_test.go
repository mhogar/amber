package integration_test

import (
	"testing"

	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserRoleCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *UserRoleCRUDTestSuite) TestCreateUserRole_WithInvalidUserRole_ReturnsError() {
	//arrange
	role := models.CreateUserRole(uuid.Nil, "", "")

	//act
	err := suite.Executor.CreateUserRole(role)

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRolesWithLesserRankByClientUID_GetsTheUserRolesWithLesserRankAndClientUIDOrderedByUsername() {
	//arrange
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 1, []byte("password")))
	user3 := suite.SaveUser(models.CreateUser("user3", 2, []byte("password")))

	client1 := suite.SaveClient(models.CreateNewClient("name1", "redirect.com", 0, "key.pem"))
	client2 := suite.SaveClient(models.CreateNewClient("name2", "redirect.com", 0, "key.pem"))

	role1 := suite.SaveUserRole(models.CreateUserRole(client1.UID, user1.Username, "role"))
	role2 := suite.SaveUserRole(models.CreateUserRole(client1.UID, user2.Username, "role"))
	suite.SaveUserRole(models.CreateUserRole(client1.UID, user3.Username, "role"))
	suite.SaveUserRole(models.CreateUserRole(client2.UID, user1.Username, "role"))

	//act
	roles, err := suite.Executor.GetUserRolesWithLesserRankByClientUID(client1.UID, 2)

	//assert
	suite.NoError(err)

	suite.Require().Len(roles, 2)
	suite.EqualValues(roles[0], role1)
	suite.EqualValues(roles[1], role2)

	//clean up
	suite.DeleteUser(user1)
	suite.DeleteUser(user2)
	suite.DeleteUser(user3)

	suite.DeleteClient(client1)
	suite.DeleteClient(client2)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleByUsernameAndClientUID_WhereUserRoleNotFound_ReturnsNilUserRole() {
	//act
	role, err := suite.Executor.GetUserRoleByClientUIDAndUsername(uuid.New(), "DNE")

	//assert
	suite.NoError(err)
	suite.Nil(role)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleByUsernameAndClientUID_GetsUserRoleWithUsernameAndClientUID() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	role := suite.SaveUserRole(models.CreateUserRole(client.UID, user.Username, "role"))

	//act
	resultRole, err := suite.Executor.GetUserRoleByClientUIDAndUsername(role.ClientUID, user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(role, resultRole)

	//clean up
	suite.DeleteUser(user)
	suite.DeleteClient(client)
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WithInvalidUserRole_ReturnsError() {
	//act
	_, err := suite.Executor.UpdateUserRole(models.CreateUserRole(uuid.Nil, "", ""))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WhereUserRoleIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Executor.UpdateUserRole(models.CreateUserRole(uuid.New(), "DNE", "role"))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_UpdatesUserRole() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	role := suite.SaveUserRole(models.CreateUserRole(client.UID, user.Username, "role"))

	//act
	role.Role = "new role"
	res, err := suite.Executor.UpdateUserRole(role)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultRole, err := suite.Executor.GetUserRoleByClientUIDAndUsername(role.ClientUID, role.Username)
	suite.NoError(err)
	suite.EqualValues(role, resultRole)

	//clean up
	suite.DeleteUser(user)
	suite.DeleteClient(client)
}

func (suite *UserRoleCRUDTestSuite) TestDeleteUserRole_WhereUserRoleNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Executor.DeleteUserRole(uuid.New(), "DNE")

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserRoleCRUDTestSuite) TestDeleteUserRole_DeletesUserRole() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	role := suite.SaveUserRole(models.CreateUserRole(client.UID, user.Username, "role"))

	//act
	res, err := suite.Executor.DeleteUserRole(role.ClientUID, role.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Executor.GetUserRoleByClientUIDAndUsername(role.ClientUID, role.Username)
	suite.NoError(err)
	suite.Nil(resultUser)

	//clean up
	suite.DeleteUser(user)
	suite.DeleteClient(client)
}

func TestUserRoleCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleCRUDTestSuite{})
}
