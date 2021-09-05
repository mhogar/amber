package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRoleControllerTestSuite struct {
	ControllerTestSuite
	UserRoleController controllers.CoreUserRoleController
}

func (suite *UserRoleControllerTestSuite) runValidateClientTestCases(validateFunc func(role *models.UserRole) common.CustomError) {
	suite.Run("EmptyUsername_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole("", "")

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be empty")
	})

	suite.Run("UsernameTooLong_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole(helpers.CreateStringOfLength(models.UserUsernameMaxLength+1), "")

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be longer", fmt.Sprint(models.UserUsernameMaxLength))
	})

	suite.Run("EmptyRole_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole("username", "")

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "role", "cannot be empty")
	})

	suite.Run("RoleTooLong_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole("username", helpers.CreateStringOfLength(models.UserRoleRoleMaxLength+1))

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "role", "cannot be longer", fmt.Sprint(models.UserRoleRoleMaxLength))
	})
}

func (suite *UserRoleControllerTestSuite) TestGetUserRoleForClient_WithErrorGettingUserRoleForClient_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserRoleForClient", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	//act
	role, cerr := suite.UserRoleController.GetUserRoleForClient(&suite.CRUDMock, uuid.New(), "username")

	//assert
	suite.Nil(role)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserRoleControllerTestSuite) TestGetUserRoleForClient_WithNoErrors_ReturnsRole() {
	//arrange
	clientUID := uuid.New()
	role := models.CreateUserRole("username", "role")

	suite.CRUDMock.On("GetUserRoleForClient", mock.Anything, mock.Anything).Return(role, nil)

	//act
	resRole, cerr := suite.UserRoleController.GetUserRoleForClient(&suite.CRUDMock, clientUID, role.Username)

	//assert
	suite.Equal(role, resRole)
	helpers.AssertNoError(&suite.Suite, cerr)

	suite.CRUDMock.AssertCalled(suite.T(), "GetUserRoleForClient", clientUID, role.Username)
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRolesForClient_ValidateUserRoleTestCases() {
	suite.runValidateClientTestCases(func(role *models.UserRole) common.CustomError {
		roles := make([]*models.UserRole, 1)
		roles[0] = role

		return suite.UserRoleController.UpdateUserRolesForClient(&suite.CRUDMock, uuid.New(), roles)
	})
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRolesForClient_WithErrorUpdatingUserRolesForClient_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("UpdateUserRolesForClient", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.UserRoleController.UpdateUserRolesForClient(&suite.CRUDMock, uuid.New(), nil)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserRoleControllerTestSuite) TestUpdateUserRolesForClient_WithNoErrors_ReturnsNoError() {
	//arrange
	clientUID := uuid.New()

	roles := make([]*models.UserRole, 1)
	roles[0] = models.CreateUserRole("username", "role")

	suite.CRUDMock.On("UpdateUserRolesForClient", mock.Anything, mock.Anything).Return(nil)

	//act
	cerr := suite.UserRoleController.UpdateUserRolesForClient(&suite.CRUDMock, clientUID, roles)

	//assert
	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUserRolesForClient", clientUID, roles)
}

func TestUserRoleControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleControllerTestSuite{})
}
