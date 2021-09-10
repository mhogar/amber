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
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPostUser_WithSessionRankLessThanUser_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)

	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     10,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)

	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)

	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithNoErrors_ReturnsUserData() {
	//arrange
	session := models.CreateNewSession("admin", 5)

	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
		Rank:     0,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	user := models.CreateUser(body.Username, body.Rank, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserDataResponse{
		Username: user.Username,
		PutUserBody: handlers.PutUserBody{
			Rank: user.Rank,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.CRUDMock, body.Username, body.Password, body.Rank)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.PutUser(nil, params, nil, &suite.CRUDMock)

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
	status, res := suite.CoreHandlers.PutUser(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPutUser_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithSessionRankLessThanNewUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 6,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "update user error"
	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUser_WithNoErrors_ReturnsUserData() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserBody{
		Rank: 1,
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)
	user := models.CreateUser(params[0].Value, body.Rank, nil)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUser(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserDataResponse{
		Username: user.Username,
		PutUserBody: handlers.PutUserBody{
			Rank: user.Rank,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUser", &suite.CRUDMock, user.Username, user.Rank)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	session := models.CreateNewSession("username", 0)

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorUpdatingUserPassword_ReturnsBadRequest() {
	//arrange
	body := handlers.PatchPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPasswordWithAuth", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorUpdatingUserPassword_ReturnsInternalServerError() {
	//arrange
	body := handlers.PatchPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPasswordWithAuth", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorDeletingAllOtherUserSessions_ReturnsBadRequest() {
	//arrange
	body := handlers.PatchPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPasswordWithAuth", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	message := "update user password error"
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorDeletingAllOtherUserSessions_ReturnsInternalServerError() {
	//arrange
	body := handlers.PatchPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPasswordWithAuth", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithNoErrors_ReturnsSuccess() {
	//arrange
	body := handlers.PatchPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	session := models.CreateNewSession("username", 0)

	suite.ControllersMock.On("UpdateUserPasswordWithAuth", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserSessions", mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.PatchPassword(req, nil, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserPasswordWithAuth", &suite.CRUDMock, session.Username, body.OldPassword, body.NewPassword)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteAllOtherUserSessions", &suite.CRUDMock, session.Username, session.Token)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	user := models.CreateUser("username", 6, nil)

	params := []httprouter.Param{
		{
			Key:   "username",
			Value: user.Username,
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "delete user error"
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUser", &suite.CRUDMock, params[0].Value)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}
