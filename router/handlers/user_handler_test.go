package handlers_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPostUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	user := models.CreateUser(body.Username, body.Rank, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserDataResponse{
		Username: user.Username,
		PutUserBody: handlers.PutUserBody{
			Rank: user.Rank,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.DataCRUDMock, body.Username, body.Password, body.Rank)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.PutUser(nil, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserHandlerTestSuite) TestPutUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPutUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "update user error"
	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	user := models.CreateUser(params[0].Value, body.Rank, nil)
	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserDataResponse{
		Username: user.Username,
		PutUserBody: handlers.PutUserBody{
			Rank: user.Rank,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUser", &suite.DataCRUDMock, params[0].Value, user.Rank)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	session := models.CreateNewSession("username", 0)

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorUpdatingUserPassword_ReturnsBadRequest() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorUpdatingUserPassword_ReturnsInternalServerError() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorDeletingAllOtherUserSessions_ReturnsBadRequest() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	message := "update user password error"
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorDeletingAllOtherUserSessions_ReturnsInternalServerError() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithNoErrors_ReturnsSuccess() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, session, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserPassword", &suite.DataCRUDMock, session.Username, body.OldPassword, body.NewPassword)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteAllOtherUserSessions", &suite.DataCRUDMock, session.Username, session.Token)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	message := "delete user error"
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUser", &suite.DataCRUDMock, params[0].Value)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}
