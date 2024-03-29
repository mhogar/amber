package controllers_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/controllers"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRoleControllerTestSuite struct {
	ControllerTestSuite
	UserRoleController controllers.CoreUserRoleController
}

func (suite *UserRoleControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()
	suite.UserRoleController = controllers.CoreUserRoleController{}
}

func (suite *UserRoleControllerTestSuite) runValidateUserRoleTestCases(validateFunc func(role *models.UserRole) common.CustomError) {
	suite.Run("EmptyRole_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole(uuid.New(), "username", "")

		//act
		cerr := validateFunc(role)

		//assert
		suite.CustomClientError(cerr, "role", "cannot be empty")
	})

	suite.Run("RoleTooLong_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole(uuid.New(), "username", helpers.CreateStringOfLength(models.UserRoleRoleMaxLength+1))

		//act
		cerr := validateFunc(role)

		//assert
		suite.CustomClientError(cerr, "role", "cannot be longer", fmt.Sprint(models.UserRoleRoleMaxLength))
	})
}

func (suite *UserRoleControllerTestSuite) TestCreateUserRole_ValidateUserRoleTestCases() {
	suite.runValidateUserRoleTestCases(func(role *models.UserRole) common.CustomError {
		return suite.UserRoleController.CreateUserRole(&suite.CRUDMock, role)
	})
}

func (suite *UserRoleControllerTestSuite) TestCreateUserRole_WithErrorGettingUserRoleByUsernameAndClientUID_ReturnsInternalError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	//act
	cerr := suite.UserRoleController.CreateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *UserRoleControllerTestSuite) TestCreateUser_WhereUserAlreadyHasRoleForClient_ReturnsClientError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(&models.UserRole{}, nil)

	//act
	cerr := suite.UserRoleController.CreateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomClientError(cerr, "user", "already has a role", "client")
}

func (suite *UserRoleControllerTestSuite) TestCreateUserRole_WithErrorCreatingUserRoleByUsernameAndClientUID_ReturnsInternalError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")

	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("CreateUserRole", mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.UserRoleController.CreateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *UserRoleControllerTestSuite) TestCreateUserRole_WithNoErrors_ReturnsNoError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")

	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("CreateUserRole", mock.Anything).Return(nil)

	//act
	cerr := suite.UserRoleController.CreateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomNoError(cerr)

	suite.CRUDMock.AssertCalled(suite.T(), "GetUserRoleByClientUIDAndUsername", role.ClientUID, role.Username)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateUserRole", role)
}

func (suite *UserRoleControllerTestSuite) TestGetUserRolesWithLesserRankByClientUID_WithErrorGettingUserRolesByClientUID_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserRolesWithLesserRankByClientUID", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	//act
	clients, cerr := suite.UserRoleController.GetUserRolesWithLesserRankByClientUID(&suite.CRUDMock, uuid.New(), 0)

	//assert
	suite.Nil(clients)
	suite.CustomInternalError(cerr)
}

func (suite *UserRoleControllerTestSuite) TestGetUserRolesWithLesserRankByClientUID_WithNoErrors_ReturnsUserRoles() {
	//arrange
	clientUID := uuid.New()
	rank := 5

	roles := []*models.UserRole{models.CreateUserRole(clientUID, "username", "role")}
	suite.CRUDMock.On("GetUserRolesWithLesserRankByClientUID", mock.Anything, mock.Anything).Return(roles, nil)

	//act
	resultRoles, cerr := suite.UserRoleController.GetUserRolesWithLesserRankByClientUID(&suite.CRUDMock, clientUID, rank)

	//assert
	suite.Equal(roles, resultRoles)
	suite.CustomNoError(cerr)

	suite.CRUDMock.AssertCalled(suite.T(), "GetUserRolesWithLesserRankByClientUID", clientUID, rank)
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRole_ValidateUserRoleTestCases() {
	suite.runValidateUserRoleTestCases(func(role *models.UserRole) common.CustomError {
		return suite.UserRoleController.UpdateUserRole(&suite.CRUDMock, role)
	})
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRole_WithErrorUpdatingUserRole_ReturnsInternalError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")
	suite.CRUDMock.On("UpdateUserRole", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserRoleController.UpdateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRole_WithFalseResultUpdatingUserRole_ReturnsClientError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")
	suite.CRUDMock.On("UpdateUserRole", mock.Anything).Return(false, nil)

	//act
	cerr := suite.UserRoleController.UpdateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomClientError(cerr, "no role found", role.Username, role.Role)
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRole_WithNoErrors_ReturnsNoError() {
	//arrange
	role := models.CreateUserRole(uuid.New(), "username", "role")
	suite.CRUDMock.On("UpdateUserRole", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserRoleController.UpdateUserRole(&suite.CRUDMock, role)

	//assert
	suite.CustomNoError(cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUserRole", role)
}

func (suite *UserRoleControllerTestSuite) TestDeleteUserRole_WithErrorDeletingUserRole_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("DeleteUserRole", mock.Anything, mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserRoleController.DeleteUserRole(&suite.CRUDMock, uuid.New(), "username")

	//assert
	suite.CustomInternalError(cerr)
}

func (suite *UserRoleControllerTestSuite) TestDeleteUserRole_WithFalseResultDeletingUserRole_ReturnsClientError() {
	//arrange
	clientUID := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteUserRole", mock.Anything, mock.Anything).Return(false, nil)

	//act
	cerr := suite.UserRoleController.DeleteUserRole(&suite.CRUDMock, clientUID, username)

	//assert
	suite.CustomClientError(cerr, "no role found", clientUID.String(), username)
}

func (suite *UserRoleControllerTestSuite) TestDeleteUserRole_WithNoErrors_ReturnsNoError() {
	//arrange
	clientUID := uuid.New()
	username := "username"

	suite.CRUDMock.On("DeleteUserRole", mock.Anything, mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserRoleController.DeleteUserRole(&suite.CRUDMock, clientUID, username)

	//assert
	suite.CustomNoError(cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUserRole", clientUID, username)
}

func TestUserRoleControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleControllerTestSuite{})
}
