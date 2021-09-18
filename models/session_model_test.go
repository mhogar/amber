package models_test

import (
	"testing"

	"authserver/models"
	"authserver/testing/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	helpers.CustomSuite
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

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, &SessionTestSuite{})
}
