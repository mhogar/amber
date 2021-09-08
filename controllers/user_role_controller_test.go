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
	suite.Run("EmptyRole_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole("username", uuid.Nil, "")

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "role", "cannot be empty")
	})

	suite.Run("RoleTooLong_ReturnsClientError", func() {
		//arrange
		role := models.CreateUserRole("username", uuid.Nil, helpers.CreateStringOfLength(models.UserRoleRoleMaxLength+1))

		//act
		cerr := validateFunc(role)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "role", "cannot be longer", fmt.Sprint(models.UserRoleRoleMaxLength))
	})
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
	roles[0] = models.CreateUserRole("username", uuid.Nil, "role")

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
