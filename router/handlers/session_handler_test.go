package handlers_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SessionHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithClientErrorCreatingSession_ReturnsBadRequest() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create session error"
	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithInternalErrorCreatingSession_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithNoErrors_ReturnsSessionData() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession(body.Username)
	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(session, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.SessionDataResponse{
		Token:    session.Token.String(),
		Username: body.Username,
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateSession", &suite.DataCRUDMock, body.Username, body.Password)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithClientErrorDeletingSession_ReturnsBadRequest() {
	//arrange
	session := &models.Session{}

	message := "delete session error"
	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithInternalErrorDeletingSession_ReturnsInternalServerError() {
	//arrange
	session := &models.Session{}

	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithNoErrors_ReturnsSuccess() {
	//arrange
	session := &models.Session{}

	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteSession", &suite.DataCRUDMock, session.Token)
}

func TestSessionHandlerTestSuite(t *testing.T) {
	suite.Run(t, &SessionHandlerTestSuite{})
}
