package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	suite.Suite
	Session *models.Session
}

func (suite *SessionTestSuite) SetupTest() {
	suite.Session = models.CreateNewSession(
		models.CreateNewUser("username", []byte("password")),
	)
}

func (suite *SessionTestSuite) TestCreateNewSession_CreatesSessionWithSuppliedFields() {
	//arrange
	user := models.CreateNewUser("", nil)

	//act
	session := models.CreateNewSession(user)

	//assert
	suite.Require().NotNil(session)
	suite.NotEqual(session.ID, uuid.Nil)
	suite.Equal(session.User, user)
}

func (suite *SessionTestSuite) TestValidate_WithValidSession_ReturnsValid() {
	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionValid, verr)
}

func (suite *SessionTestSuite) TestValidate_WithNilID_ReturnsSessionInvalidID() {
	//arrange
	suite.Session.ID = uuid.Nil

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionNilID, verr)
}

func (suite *SessionTestSuite) TestValidate_WithNilUser_ReturnsSessionNilUser() {
	//arrange
	suite.Session.User = nil

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionNilUser, verr)
}

func (suite *SessionTestSuite) TestValidate_WithInvalidUser_ReturnsSessionInvalidUser() {
	//arrange
	suite.Session.User = models.CreateNewUser("", nil)

	//act
	verr := suite.Session.Validate()

	//assert
	suite.Equal(models.ValidateSessionInvalidUser, verr)
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, &SessionTestSuite{})
}
