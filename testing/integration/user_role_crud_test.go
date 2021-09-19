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
	err := suite.Tx.CreateUserRole(role)

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
	roles, err := suite.Tx.GetUserRolesWithLesserRankByClientUID(client1.UID, 2)

	//assert
	suite.NoError(err)

	suite.Require().Len(roles, 2)
	suite.EqualValues(roles[0], role1)
	suite.EqualValues(roles[1], role2)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleByUsernameAndClientUID_WhereUserRoleNotFound_ReturnsNilUserRole() {
	//act
	role, err := suite.Tx.GetUserRoleByClientUIDAndUsername(uuid.New(), "DNE")

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
	resultRole, err := suite.Tx.GetUserRoleByClientUIDAndUsername(role.ClientUID, user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(role, resultRole)
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WithInvalidUserRole_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUserRole(models.CreateUserRole(uuid.Nil, "", ""))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WhereUserRoleIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUserRole(models.CreateUserRole(uuid.New(), "DNE", "role"))

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
	res, err := suite.Tx.UpdateUserRole(role)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultRole, err := suite.Tx.GetUserRoleByClientUIDAndUsername(role.ClientUID, role.Username)
	suite.NoError(err)
	suite.EqualValues(role, resultRole)
}

func (suite *UserRoleCRUDTestSuite) TestDeleteUserRole_WhereUserRoleNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.DeleteUserRole("DNE", uuid.New())

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
	res, err := suite.Tx.DeleteUserRole(role.Username, role.ClientUID)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserRoleByClientUIDAndUsername(role.ClientUID, role.Username)
	suite.NoError(err)
	suite.Nil(resultUser)
}

func TestUserRoleCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleCRUDTestSuite{})
}
