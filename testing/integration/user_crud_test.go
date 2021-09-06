package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *UserCRUDTestSuite) TestCreateUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.CreateUser(models.CreateUser("", nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestGetUserByUsername_WhereUserNotFound_ReturnsNilUser() {
	//act
	user, err := suite.Tx.GetUserByUsername("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(user)
}

func (suite *UserCRUDTestSuite) TestGetUserByUsername_GetsTheUserWithUsername() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", []byte("password")))

	//act
	resultUser, err := suite.Tx.GetUserByUsername(user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithInvalidUser_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUser(models.CreateUser("", nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUser(models.CreateUser("username", []byte("password")))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_UpdatesUserWithId() {
	//arrange
	newPassword := []byte("new_password")
	user := suite.SaveUser(models.CreateUser("username", []byte("password")))

	//act
	user.PasswordHash = newPassword
	res, err := suite.Tx.UpdateUser(user)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.DeleteUser("not_a_real_username")

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_DeletesUserWithId() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", []byte("password")))

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.Nil(resultUser)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_AlsoDeletesAllRolesForUser() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	suite.UpdateUserRolesForClient(client, models.CreateUserRole(user.Username, "role"))

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	roles, err := suite.Tx.GetUserRoleForClient(client.UID, user.Username)
	suite.NoError(err)
	suite.Empty(roles)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_AlsoDeletesAllUserSessions() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", []byte("password")))
	session := suite.SaveSession(models.CreateNewSession(user.Username))

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultSession, err := suite.Tx.GetSessionByToken(session.Token)
	suite.NoError(err)
	suite.Nil(resultSession)
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserCRUDTestSuite{})
}
