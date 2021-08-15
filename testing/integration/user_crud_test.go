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
	err := suite.Tx.CreateUser(models.CreateNewUser("", nil))

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

func (suite *UserCRUDTestSuite) TestGetUserByUsernameGetsTheUserWithUsername() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	suite.CreateUser(user)

	//act
	resultUser, err := suite.Tx.GetUserByUsername(user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithInvalidUser_ReturnsError() {
	//act
	_, err := suite.Tx.UpdateUser(models.CreateNewUser("", nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WhereUserIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.UpdateUser(models.CreateNewUser("username", []byte("password")))

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_UpdatesUserWithId() {
	//arrange
	newPassword := []byte("new_password")

	user := models.CreateNewUser("username", []byte("password"))
	suite.CreateUser(user)

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
	user := models.CreateNewUser("username", []byte("password"))
	suite.CreateUser(user)

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByUsername(user.Username)
	suite.NoError(err)
	suite.Nil(resultUser)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_AlsoDeletesAllUserTokens() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	token := models.CreateNewAccessToken(
		user,
		models.CreateNewClient("name"),
	)
	suite.SaveAccessTokenAndFields(token)

	//act
	res, err := suite.Tx.DeleteUser(user.Username)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)
	suite.NoError(err)
	suite.Nil(resultAccessToken)
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserCRUDTestSuite{})
}
