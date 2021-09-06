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

func (suite *UserRoleHandlerTestSuite) TestPutClientRoles_WithErrorParsingId_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, nil)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClientRoles(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "client id", "invalid format")
}

func (suite *UserRoleHandlerTestSuite) TestPutClientRoles_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClientRoles(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *UserRoleHandlerTestSuite) TestPutClientRoles_WithClientErrorUpdatingClientRoles_ReturnsBadRequest() {
	//arrange
	rolesBody := make([]handlers.PutClientRolesBody, 1)
	rolesBody[0] = handlers.PutClientRolesBody{
		Username: "username",
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, rolesBody)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	message := "update client roles error"
	suite.ControllersMock.On("UpdateUserRolesForClient", mock.Anything, mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutClientRoles(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *UserRoleHandlerTestSuite) TestPutClientRoles_WithInternalErrorUpdatingClientRoles_ReturnsInternalServerError() {
	//arrange
	rolesBody := make([]handlers.PutClientRolesBody, 1)
	rolesBody[0] = handlers.PutClientRolesBody{
		Username: "username",
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, rolesBody)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("UpdateUserRolesForClient", mock.Anything, mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutClientRoles(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *UserRoleHandlerTestSuite) TestPutClientRoles_WithNoErrors_ReturnsInternalServerError() {
	//arrange
	rolesBody := make([]handlers.PutClientRolesBody, 1)
	rolesBody[0] = handlers.PutClientRolesBody{
		Username: "username",
		Role:     "role",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, rolesBody)

	clientUID := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: clientUID.String(),
		},
	}

	suite.ControllersMock.On("UpdateUserRolesForClient", mock.Anything, mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.PutClientRoles(req, params, nil, &suite.DataCRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateUserRolesForClient", &suite.DataCRUDMock, clientUID, mock.MatchedBy(func(roles []*models.UserRole) bool {
		return len(rolesBody) == len(roles) &&
			rolesBody[0].Username == roles[0].Username &&
			rolesBody[0].Role == roles[0].Role
	}))
}

func TestUserRoleHandlerTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleHandlerTestSuite{})
}
