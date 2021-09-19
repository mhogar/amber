package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/controllers/mocks"
	"authserver/models"
	"errors"
	"testing"

	"github.com/google/uuid"
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
	authErr := common.ClientError("authenticate user error")
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, authErr)

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, "username", "password")

	//assert
	suite.Nil(session)
	suite.Equal(cerr, authErr)
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithErrorSavingSession_ReturnsInternalError() {
	//arrange
	user := models.CreateUser("username", 0, nil)

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(errors.New(""))

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, user.Username, "password")

	//assert
	suite.Nil(session)
	suite.CustomInternalError(cerr)
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithNoErrors_ReturnsNoError() {
	//arrange
	user := models.CreateUser("username", 0, nil)
	password := "password"

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(nil)

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, user.Username, password)

	//assert
	suite.Require().NotNil(session)
	suite.Equal(user.Username, session.Username)
	suite.Equal(user.Rank, session.Rank)
	suite.CustomNoError(cerr)

	suite.ControllersMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, user.Username, password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveSession", session)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithErrorDeletingSession_ReturnsInternalError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithFalseResultDeletingSession_ReturnsClientError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(false, nil)

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	suite.CustomClientError(cerr, "session with id", id.String(), "not found")
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithNoErrors_ReturnsNoError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(true, nil)

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	suite.CustomNoError(cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteSession", id)
}

func (suite *SessionControllerTestSuite) TestDeleteAllUserSessions_WithErrorDeletingSessions_ReturnsInternalError() {
	//arrange
	username := "username"

	suite.CRUDMock.On("DeleteAllUserSessions", mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.SessionController.DeleteAllUserSessions(&suite.CRUDMock, username)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllUserSessions_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"

	suite.CRUDMock.On("DeleteAllUserSessions", mock.Anything).Return(nil)

	//act
	cerr := suite.SessionController.DeleteAllUserSessions(&suite.CRUDMock, username)

	//assert
	suite.CustomNoError(cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllUserSessions", username)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithErrorDeletingSessions_ReturnsInternalError() {
	//arrange
	id := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, username, id)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithNoErrors_ReturnsNoError() {
	//arrange
	id := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything).Return(nil)

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, username, id)

	//assert
	suite.CustomNoError(cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllOtherUserSessions", username, id)
}

func TestSessionControllerTestSuite(t *testing.T) {
	suite.Run(t, &SessionControllerTestSuite{})
}
