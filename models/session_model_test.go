package models_test

import (
	"testing"

	"authserver/models"
	"authserver/testing/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	suite.Suite
	Session *models.Session
}

func (suite *SessionTestSuite) SetupTest() {
	suite.Session = models.CreateNewSession("username", 0)
}

func (suite *SessionTestSuite) TestCreateNewSession_CreatesSessionWithSuppliedFields() {
	//arrange
	username := "username"
	rank := 100

	//act
	session := models.CreateNewSession(username, rank)

	//assert
	suite.Require().NotNil(session)
	suite.NotEqual(session.Token, uuid.Nil)
	suite.Equal(session.Username, username)
	suite.Equal(session.Rank, rank)
}

func (suite *SessionTestSuite) TestValidate_WithValidSession_ReturnsValid() {
	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionValid, verr)
}

func (suite *SessionTestSuite) TestValidate_WithNilID_ReturnsSessionNilToken() {
	//arrange
	suite.Session.Token = uuid.Nil

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionNilToken, verr)
}

func (suite *SessionTestSuite) TestValidate_WithEmptyUsername_ReturnsSessionEmptyUsername() {
	//arrange
	suite.Session.Username = ""

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionEmptyUsername, verr)
}

func (suite *SessionTestSuite) TestValidate_UsernameMaxLengthTestCases() {
	var username string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.Session.Username = username

		//act
		verr := suite.Session.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	username = helpers.CreateStringOfLength(models.UserUsernameMaxLength)
	expectedValidateError = models.ValidateSessionValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	username += "a"
	expectedValidateError = models.ValidateSessionUsernameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func (suite *SessionTestSuite) TestValidate_WithNegativeRank_ReturnsSessionInvalidRank() {
	//arrange
	suite.Session.Rank = -1

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionInvalidRank, verr)
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, &SessionTestSuite{})
}
