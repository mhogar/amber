package handlers_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router/handlers"
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

func (suite *UserRoleHandlerTestSuite) TestGetUserRoles_WithInvalidClientID_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.GetUserRoles(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestGetUserRoles_WithClientErrorGettingUserRolesWithLesserRankByClientUID_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}
	session := models.CreateNewSession("admin", 5)

	message := "get user-roles error"
	suite.ControllersMock.On("GetUserRolesWithLesserRankByClientUID", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.GetUserRoles(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestGetUserRoles_WithInternalErrorGettingUserRolesWithLesserRankByClientUID_ReturnsInternalServerError() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}
	session := models.CreateNewSession("admin", 5)

	suite.ControllersMock.On("GetUserRolesWithLesserRankByClientUID", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.GetUserRoles(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestGetUserRoles_WithNoErrors_ReturnsUserRoleData() {
	//arrange
	clientUID := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: clientUID.String(),
		},
	}
	session := models.CreateNewSession("admin", 5)

	roles := []*models.UserRole{
		models.CreateUserRole(clientUID, "user1", "role"),
		models.CreateUserRole(clientUID, "user2", "role"),
	}
	suite.ControllersMock.On("GetUserRolesWithLesserRankByClientUID", mock.Anything, mock.Anything, mock.Anything).Return(roles, common.NoError())

	//act
	status, res := suite.CoreHandlers.GetUserRoles(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, []handlers.UserRoleDataResponse{
		{
			PostUserRoleBody: handlers.PostUserRoleBody{
				Username: roles[0].Username,
				Role:     roles[0].Role,
			},
		},
		{
			PostUserRoleBody: handlers.PostUserRoleBody{
				Username: roles[1].Username,
				Role:     roles[1].Role,
			},
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "GetUserRolesWithLesserRankByClientUID", &suite.CRUDMock, clientUID, session.Rank)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithInvalidClientID_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.PostUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}
	req := suite.CreateDummyJSONRequest("invalid")

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "invalid json body")
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	suite.InsufficientPermissionsErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithClientErrorCreatingUserRole_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "create user role error"
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithInternalErrorCreatingUserRole_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPostUserRole_WithNoErrors_ReturnsUserRoleData() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	body := handlers.PostUserRoleBody{
		Username: "username",
		Role:     "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	var role *models.UserRole
	suite.ControllersMock.On("CreateUserRole", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		role = args.Get(1).(*models.UserRole)
	})

	//act
	status, res := suite.CoreHandlers.PostUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, handlers.UserRoleDataResponse{
		PostUserRoleBody: handlers.PostUserRoleBody{
			Username: role.Username,
			Role:     role.Role,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, body.Username, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUserRole", &suite.CRUDMock, role)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithErrorParsingClientId_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}
	req := suite.CreateDummyJSONRequest(nil)

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	//act
	status, res := suite.CoreHandlers.PutUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "username not provided")
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}
	req := suite.CreateDummyJSONRequest("invalid")

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "invalid json body")
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{

		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{

		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	suite.InsufficientPermissionsErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithClientErrorUpdatingUserRole_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "create user role error"
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithInternalErrorUpdatingUserRole_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestPutUserRole_WithNoErrors_ReturnsUserData() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	body := handlers.PutUserRoleBody{
		Role: "role",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	var role *models.UserRole
	suite.ControllersMock.On("UpdateUserRole", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		role = args.Get(1).(*models.UserRole)
	})

	//act
	status, res := suite.CoreHandlers.PutUserRole(req, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, handlers.UserRoleDataResponse{
		PostUserRoleBody: handlers.PostUserRoleBody{
			Username: role.Username,
			Role:     role.Role,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[1].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserRole", &suite.CRUDMock, role)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithErrorParsingClientId_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}
	req := suite.CreateDummyJSONRequest(nil)

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithMissingUsername_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "username not provided")
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithClientErrorVerifyingUserRank_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	message := "verify user rank error"
	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithInternalErrorVerifyingUserRank_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithFalseResultVerifyingUserRank_ReturnsForbidden() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	user := models.CreateUser("username", 6, nil)

	params := []httprouter.Param{
		{
			Key:   "username",
			Value: user.Username,
		},
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(false, common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusForbidden, status)
	suite.InsufficientPermissionsErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())

	message := "delete user error"
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
		{
			Key:   "username",
			Value: "username",
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *UserRoleHandlerTestSuite) TestDeleteUserRole_WithNoErrors_ReturnsSuccess() {
	//arrange
	session := models.CreateNewSession("admin", 5)
	clientID := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "username",
			Value: "username",
		},
		{
			Key:   "id",
			Value: clientID.String(),
		},
	}

	suite.ControllersMock.On("VerifyUserRank", mock.Anything, mock.Anything, mock.Anything).Return(true, common.NoError())
	suite.ControllersMock.On("DeleteUserRole", mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteUserRole(nil, params, session, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessResponse(res)

	suite.ControllersMock.AssertCalled(suite.T(), "VerifyUserRank", &suite.CRUDMock, params[0].Value, session.Rank)
	suite.ControllersMock.AssertCalled(suite.T(), "DeleteUserRole", &suite.CRUDMock, params[0].Value, clientID)
}

func TestUserRoleHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleHandlerTestSuite{})
}
