package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/controllers/mocks"
	"authserver/testing/helpers"
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
	username := "username"
	password := "password"

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(errors.New(""))

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(session)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestCreateSession_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"
	password := "password"

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("SaveSession", mock.Anything).Return(nil)

	//act
	session, cerr := suite.SessionController.CreateSession(&suite.CRUDMock, username, password)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, username, password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveSession", session)

	suite.Require().NotNil(session)
	suite.Equal(username, session.Username)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithErrorDeletingSession_ReturnsInternalError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithFalseResultDeletingSession_ReturnsClientError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(false, nil)

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "session with id", id.String(), "not found")
}

func (suite *SessionControllerTestSuite) TestDeleteSession_WithNoErrors_ReturnsNoError() {
	//arrange
	id := uuid.New()
	suite.CRUDMock.On("DeleteSession", mock.Anything).Return(true, nil)

	//act
	cerr := suite.SessionController.DeleteSession(&suite.CRUDMock, id)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteSession", id)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithErrorDeletingSessions_ReturnsInternalError() {
	//arrange
	id := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, username, id)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *SessionControllerTestSuite) TestDeleteAllOtherUserSessions_WithNoErrors_ReturnsNoError() {
	//arrange
	id := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything).Return(nil)

	//act
	cerr := suite.SessionController.DeleteAllOtherUserSessions(&suite.CRUDMock, username, id)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllOtherUserSessions", username, id)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func TestSessionControlTestSuite(t *testing.T) {
	suite.Run(t, &SessionControllerTestSuite{})
}
