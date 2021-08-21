package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/controllers/mocks"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
)

type SessionControllerTestSuite struct {
	ControllerTestSuite
	ControllersMock   mocks.Controllers
	SessionController controllers.CoreSessionController
}

func (suite *SessionControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.ControllersMock = mocks.Controllers{}
	suite.SessionController = controllers.CoreSessionController{
		AuthController: &suite.ControllersMock,
	}
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithErrorAuthenticatingUser_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	authErr := common.ClientError("authenticate user error")
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, authErr)

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(session)
	suite.Equal(cerr, authErr)
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithErrorSavingSession_ReturnsInternalError() {
	//arrange
	password := "password"
	user := models.CreateUser("username", []byte(password))

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(errors.New(""))

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, user.Username, password)

	//assert
	suite.Nil(session)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithNoErrors_ReturnsNoError() {
	//arrange
	password := "password"
	user := models.CreateUser("username", []byte(password))

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(nil)

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, user.Username, password)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, user.Username, password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveSession", session)

	suite.Require().NotNil(session)
	suite.Equal(user, session.User)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithErrorDeletingSession_ReturnsInternalError() {
	//arrange
	session := &models.Session{}
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, session)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithNoErrors_ReturnsNoError() {
	//arrange
	session := &models.Session{}
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(nil)

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, session)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteSession", session)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithErrorDeletingSessions_ReturnsInternalError() {
	//arrange
	session := &models.Session{}
	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, session)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithNoErrors_ReturnsNoError() {
	//arrange
	session := &models.Session{}
	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything).Return(nil)

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, session)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllOtherUserSessions", session)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func TestSessionControlTestSuite(t *testing.T) {
	suite.Run(t, &SessionControllerTestSuite{})
}
