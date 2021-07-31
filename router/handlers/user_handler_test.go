package handlers_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/models"
	"authserver/router"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	RouterTestSuite
}

func (suite *UserHandlerTestSuite) TestPostUser_WithErrorCreatingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", nil)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", "invalid")

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPostUser_WithClientErrorCreatingUser_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	message := "create user error"
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	//actFSuc
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithInternalErrorCreatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithErrorCommitingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertNotCalled(suite.T(), "Authenticate", mock.Anything)
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.TransactionMock, body.Username, body.Password)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUser_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PostUserBody{
		Username: "username",
		Password: "password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPost, server.URL+"/user", "", body)

	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithErrorCreatingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	message := "delete user error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithInternalDeletingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithErrorCommitingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUser", &suite.TransactionMock, token.User)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodDelete, server.URL+"/user", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("DeleteUser", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorAuthenticatingUser_ReturnsUnauthorized() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", nil)

	message := "authenticate error"
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertErrorResponse(&suite.Suite, res, http.StatusUnauthorized, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorAuthenticatingUser_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(nil, requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithErrorCreatingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	token := &models.AccessToken{User: &models.User{}}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", nil)

	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", "invalid")

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorUpdatingUserPassword_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorUpdatingUserPassword_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithClientErrorDeletingAllOtherUserTokens_ReturnsBadRequest() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	message := "update user password error"
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(requesterror.ClientError(message))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, message)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithInternalErrorDeletingAllOtherUserTokens_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(requesterror.InternalError())

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithErrorCommitingTransaction_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(""))

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithValidRequest_ReturnsSuccess() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.ControllersMock.On("DeleteAllOtherUserTokens", mock.Anything, mock.Anything).Return(requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.AuthenticatorMock.AssertCalled(suite.T(), "Authenticate", mock.Anything)
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserPassword", &suite.TransactionMock, token.User, body.OldPassword, body.NewPassword)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteAllOtherUserTokens", &suite.TransactionMock, token)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	common.AssertSuccessResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestUpdateUserPassword_WithPanicTriggered_ReturnsInternalServerError() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	body := router.PatchUserPasswordBody{
		OldPassword: "old password",
		NewPassword: "new password",
	}
	req := common.CreateRequest(&suite.Suite, http.MethodPatch, server.URL+"/user/password", "", body)

	token := &models.AccessToken{User: &models.User{}}
	suite.AuthenticatorMock.On("Authenticate", mock.Anything).Return(token, requesterror.NoError())
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("UpdateUserPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	common.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserHandlerTestSuite{})
}
