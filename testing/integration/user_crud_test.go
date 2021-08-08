package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *UserCRUDTestSuite) TestSaveUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.SaveUser(models.CreateNewUser("", nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestGetUserById_WhereUserNotFound_ReturnsNilUser() {
	//act
	user, err := suite.Tx.GetUserByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(user)
}

func (suite *UserCRUDTestSuite) TestGetUserById_GetsTheUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	suite.SaveUser(user)

	//act
	resultUser, err := suite.Tx.GetUserByID(user.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
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
	suite.SaveUser(user)

	//act
	resultUser, err := suite.Tx.GetUserByUsername(user.Username)

	//assert
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithInvalidUser_ReturnsError() {
	//act
	err := suite.Tx.UpdateUser(models.CreateNewUser("", nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "user model")
}

func (suite *UserCRUDTestSuite) TestUpdateUser_WithNoUserToUpdate_ReturnsNilError() {
	//act
	err := suite.Tx.UpdateUser(models.CreateNewUser("username", []byte("password")))

	//assert
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestUpdateUser_UpdatesUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	suite.SaveUser(user)

	//act
	user.Username = "username2"
	err := suite.Tx.UpdateUser(user)

	//assert
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByID(user.ID)
	suite.NoError(err)
	suite.EqualValues(user, resultUser)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_WithNoUserToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteUser(models.CreateNewUser("", nil))

	//assert
	suite.NoError(err)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_DeletesUserWithId() {
	//arrange
	user := models.CreateNewUser("username", []byte("password"))
	suite.SaveUser(user)

	//act
	err := suite.Tx.DeleteUser(user)

	//assert
	suite.Require().NoError(err)

	resultUser, err := suite.Tx.GetUserByID(user.ID)
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
	err := suite.Tx.DeleteUser(user)

	//assert
	suite.Require().NoError(err)

	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)
	suite.NoError(err)
	suite.Nil(resultAccessToken)
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, &UserCRUDTestSuite{})
}
