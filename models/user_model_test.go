package models_test

import (
	"testing"

	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	helpers.CustomSuite
	User *models.User
}

func (suite *UserTestSuite) SetupTest() {
	suite.User = models.CreateUser("username", 0, []byte("password"))
}

func (suite *UserTestSuite) TestCreateNewUser_CreatesUserWithSuppliedFields() {
	//arrange
	username := "this is a test username"
	rank := 100
	hash := []byte("this is a password")

	//act
	user := models.CreateUser(username, rank, hash)

	//assert
	suite.Require().NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(hash, user.PasswordHash)
	suite.Equal(rank, user.Rank)
}

func (suite *UserTestSuite) TestValidate_WithValidUser_ReturnsValid() {
	//act
	verr := suite.User.Validate()

	//assert
	suite.Equal(models.ValidateUserValid, verr)
}

func (suite *UserTestSuite) TestValidate_WithEmptyUsername_ReturnsUserEmptyUsername() {
	//arrange
	suite.User.Username = ""

	//act
	verr := suite.User.Validate()

	//assert
	suite.Equal(models.ValidateUserEmptyUsername, verr)
}

func (suite *UserTestSuite) TestValidate_UsernameMaxLengthTestCases() {
	var username string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.User.Username = username

		//act
		verr := suite.User.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	username = helpers.CreateStringOfLength(models.UserUsernameMaxLength)
	expectedValidateError = models.ValidateUserValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	username += "a"
	expectedValidateError = models.ValidateUserUsernameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func (suite *UserTestSuite) TestValidate_WithNegativeRank_ReturnsUserInvalidRank() {
	//arrange
	suite.User.Rank = -1

	//act
	verr := suite.User.Validate()

	//assert
	suite.Equal(models.ValidateUserInvalidRank, verr)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, &UserTestSuite{})
}
