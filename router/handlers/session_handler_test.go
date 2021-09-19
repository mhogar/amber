package handlers_test

import (
	"net/http"
	"testing"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/handlers"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SessionHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithInvalidJSONBody_ReturnsInvalidRequest() {
	//arrange
	req := suite.CreateDummyJSONRequest("invalid")

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "invalid json body")
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithClientErrorCreatingSession_ReturnsBadRequest() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := suite.CreateDummyJSONRequest(body)

	message := "create session error"
	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithInternalErrorCreatingSession_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *SessionHandlerTestSuite) TestPostSession_WithNoErrors_ReturnsSessionData() {
	//arrange
	body := handlers.PostSessionBody{
		Username: "username",
		Password: "password",
	}
	req := suite.CreateDummyJSONRequest(body)

	session := models.CreateNewSession(body.Username, 0)
	suite.ControllersMock.On("CreateSession", mock.Anything, mock.Anything, mock.Anything).Return(session, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostSession(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, handlers.SessionDataResponse{
		Token:    session.Token.String(),
		Username: body.Username,
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateSession", &suite.CRUDMock, body.Username, body.Password)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithClientErrorDeletingSession_ReturnsBadRequest() {
	//arrange
	session := &models.Session{}

	message := "delete session error"
	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithInternalErrorDeletingSession_ReturnsInternalServerError() {
	//arrange
	session := &models.Session{}

	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *SessionHandlerTestSuite) TestDeleteSession_WithNoErrors_ReturnsSuccess() {
	//arrange
	session := &models.Session{}

	suite.ControllersMock.On("DeleteSession", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteSession(nil, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessResponse(res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteSession", &suite.CRUDMock, session.Token)
}

func TestSessionHandlerTestSuite(t *testing.T) {
	suite.Run(t, &SessionHandlerTestSuite{})
}
