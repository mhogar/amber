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

func (suite *UserRoleCRUDTestSuite) TestCreateUserRole_WithInvalidUserRole_ReturnsError() {
	//arrange
	role := models.CreateUserRole("", uuid.Nil, "")

	//act
	err := suite.Tx.CreateUserRole(role)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRolesByClientUID_GetsUserRolesWithClientUIDOrderedByUsername() {
	//arrange
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 0, []byte("password")))

	client1 := suite.SaveClient(models.CreateNewClient("name1", "redirect.com", 0, "key.pem"))
	client2 := suite.SaveClient(models.CreateNewClient("name2", "redirect.com", 0, "key.pem"))

	role1 := suite.SaveUserRole(models.CreateUserRole(user1.Username, client1.UID, "role"))
	role2 := suite.SaveUserRole(models.CreateUserRole(user2.Username, client1.UID, "role"))
	suite.SaveUserRole(models.CreateUserRole(user1.Username, client2.UID, "role"))

	//act
	roles, err := suite.Tx.GetUserRolesByClientUID(client1.UID)

	//assert
	suite.NoError(err)

	suite.Require().Len(roles, 2)
	suite.EqualValues(roles[0], role1)
	suite.EqualValues(roles[1], role2)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleByUsernameAndClientUID_WhereUserRoleNotFound_ReturnsNilUserRole() {
	//act
	role, err := suite.Tx.GetUserRoleByUsernameAndClientUID("DNE", uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(role)
}

func (suite *UserRoleCRUDTestSuite) TestGetUserRoleByUsernameAndClientUID_GetsUserRoleWithUsernameAndClientUID() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	role := suite.SaveUserRole(models.CreateUserRole(user.Username, client.UID, "role"))

	//act
	resultRole, err := suite.Tx.GetUserRoleByUsernameAndClientUID(user.Username, role.ClientUID)

	//assert
	suite.NoError(err)
	suite.EqualValues(role, resultRole)
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WithInvalidUserRole_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUserRole(models.CreateUserRole("", uuid.Nil, ""))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user-role model")
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_WhereUserRoleIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUserRole(models.CreateUserRole("DNE", uuid.New(), "role"))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserRoleCRUDTestSuite) TestUpdateUserRole_UpdatesUserRole() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	role := suite.SaveUserRole(models.CreateUserRole(user.Username, client.UID, "role"))

	//act
	role.Role = "new role"
	res, err := suite.Tx.UpdateUserRole(role)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultRole, err := suite.Tx.GetUserRoleByUsernameAndClientUID(role.Username, role.ClientUID)
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
	role := suite.SaveUserRole(models.CreateUserRole(user.Username, client.UID, "role"))

	//act
	res, err := suite.Tx.DeleteUserRole(role.Username, role.ClientUID)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserRoleByUsernameAndClientUID(role.Username, role.ClientUID)
	suite.NoError(err)
	suite.Nil(resultUser)
}

func TestUserRoleCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleCRUDTestSuite{})
}
