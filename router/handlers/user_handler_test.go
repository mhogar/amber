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

type UserHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPostUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	body := handlers.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostUser(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.TransactionMock, body.Username, body.Password)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	token := &models.AccessToken{User: &models.User{}}

	message := "delete user error"
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	token := &models.AccessToken{User: &models.User{}}

	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithNoErrors_ReturnsSuccess() {
	//arrange
	token := &models.AccessToken{User: &models.User{}}

	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUser(nil, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUser", &suite.TransactionMock, token.User.Username)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	token := &models.AccessToken{User: &models.User{}}

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

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

	token := &models.AccessToken{User: &models.User{}}

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

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

	token := &models.AccessToken{User: &models.User{}}

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorDeletingAllOtherUserTokens_ReturnsBadRequest() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	token := &models.AccessToken{User: &models.User{}}

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorDeletingAllOtherUserTokens_ReturnsInternalServerError() {
	//arrange
	body := handlers.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	token := &models.AccessToken{User: &models.User{}}

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

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

	token := &models.AccessToken{User: &models.User{}}

	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.PatchUserPassword(req, nil, token, &suite.TransactionMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserPassword", &suite.TransactionMock, token.User, body.OldPassword, body.NewPassword)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteAllOtherUserTokens", &suite.TransactionMock, token)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}
