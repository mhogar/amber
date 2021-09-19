package integration_test

import (
	"testing"

	"github.com/mhogar/amber/models"

	"github.com/stretchr/testify/suite"
)

type UserCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *UserCRUDTestSuite) TestCreateUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.CreateUser(models.CreateUser("", -1, nil))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestCreateUser_WithNilPasswordHash_ReturnsError() {
	//act
	err := suite.Tx.CreateUser(models.CreateUser("username", 0, nil))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "password hash", "cannot be nil")
}

func (suite *UserCRUDTestSuite) TestGetUsersWithLesserRank_GetsTheUsersWithLesserRankOrderedByUsername() {
	//arrange
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 1, []byte("password")))
	suite.SaveUser(models.CreateUser("user3", 2, []byte("password")))

	//act
	users, err := suite.Tx.GetUsersWithLesserRank(2)

	//assert
	suite.NoError(err)

	suite.Require().Len(users, 2)
	suite.EqualValues(users[0], user1)
	suite.EqualValues(users[1], user2)
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
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))

	//act
	resultUser, err := suite.Tx.GetUserByUsername(user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithInvalidUser_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUser(models.CreateUser("", -1, nil))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUser(models.CreateUser("DNE", 0, []byte("password")))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_UpdatesUser() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))

	//act
	user.Rank = 10
	res, err := suite.Tx.UpdateUser(user)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUserPassword_WithNilHash_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUserPassword("username", nil)

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "password hash", "cannot be nil")
}

func (suite *UserCRUDTestSuite) TestUpdateUserPassword_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUserPassword("username", []byte("password"))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUserPassword_UpdatesUserWithUsername() {
	//arrange
	newPassword := []byte("new_password")
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))

	//act
	res, err := suite.Tx.UpdateUserPassword(user.Username, newPassword)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.Equal(newPassword, resultUser.PasswordHash)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.DeleteUser("not_a_real_username")

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_DeletesUserWithUsername() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))

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
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	suite.SaveUserRole(models.CreateUserRole(client.UID, user.Username, "role"))

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	role, err := suite.Tx.GetUserRoleByClientUIDAndUsername(client.UID, user.Username)
	suite.NoError(err)
	suite.Nil(role)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_AlsoDeletesAllUserSessions() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	session := suite.SaveSession(models.CreateNewSession(user.Username, 0))

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
