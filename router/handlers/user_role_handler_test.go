package handlers_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRoleHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.PostUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithClientErrorCreatingUserRole_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "create user role error"
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithInternalErrorCreatingUserRole_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPostUserRole_WithNoErrors_ReturnsUserRoleData() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PostUserRoleBody{
		ClientID: uuid.New(),
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	var role *models.UserRole
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		role = args.Get(1).(*models.UserRole)
	})

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserRoleDataResponse{
		Username: role.Username,
		PostUserRoleBody: handlers.PostUserRoleBody{
			ClientID: role.ClientUID,
			Role:     role.Role,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUserRole", &suite.CRUDMock, role)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.PutUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithErrorParsingClientId_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: "invalid",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, nil)

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithClientErrorUpdatingUserRole_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "create user role error"
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithInternalErrorUpdatingUserRole_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestPutUserRole_WithNoErrors_ReturnsUserData() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	var role *models.UserRole
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		role = args.Get(1).(*models.UserRole)
	})

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.UserRoleDataResponse{
		Username: role.Username,
		PostUserRoleBody: handlers.PostUserRoleBody{
			ClientID: role.ClientUID,
			Role:     role.Role,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserRole", &suite.CRUDMock, role)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{}

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "username not provided")
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithErrorParsingClientId_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: "invalid",
		},
	}
	req := helpers.CreateDummyRequest(&suite.Suite, nil)

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "client id", "invalid format")
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	user := models.CreateUser("username", 6, nil)

	params := []httprouter.Param{
		{
			Key:   "username",
			Value: user.Username,
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	helpers.AssertInsufficientPermissionsErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "delete user error"
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserHandlerTestSuite) TestDeleteUserRole_WithNoErrors_ReturnsSuccess() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	clientID := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "client_id",
			Value: clientID.String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUserRole", &suite.CRUDMock, params[0].Value, clientID)
}

func TestUserRoleHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleHandlerTestSuite{})
}
