package models_test

import (
	"testing"

	"authserver/models"
	"authserver/testing/helpers"

	"github.com/stretchr/testify/suite"
)

type UserRoleTestSuite struct {
	suite.Suite
	UserRole *models.UserRole
}

func (suite *UserRoleTestSuite) SetupTest() {
	suite.UserRole = models.CreateUserRole("username", "role")
}

func (suite *UserRoleTestSuite) TestCreateNewUserRole_CreatesUserRoleWithSuppliedFields() {
	//arrange
	username := "this is a test username"
	role := "this is a test role"

	//act
	userRole := models.CreateUserRole(username, role)

	//assert
	suite.Require().NotNil(userRole)
	suite.Equal(username, userRole.Username)
	suite.Equal(role, userRole.Role)
}

func (suite *UserRoleTestSuite) TestValidate_WithValidUserRole_ReturnsValid() {
	//act
	verr := suite.UserRole.Validate()

	//assert
	suite.Equal(models.ValidateUserValid, verr)
}

func (suite *UserRoleTestSuite) TestValidate_WithEmptyUsername_ReturnsUserRoleEmptyUsername() {
	//arrange
	suite.UserRole.Username = ""

	//act
	verr := suite.UserRole.Validate()

	//assert
	suite.Equal(models.ValidateUserRoleEmptyUsername, verr)
}

func (suite *UserRoleTestSuite) TestValidate_UsernameMaxLengthTestCases() {
	var username string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.UserRole.Username = username

		//act
		verr := suite.UserRole.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	username = helpers.CreateStringOfLength(models.UserUsernameMaxLength)
	expectedValidateError = models.ValidateUserRoleValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	username += "a"
	expectedValidateError = models.ValidateUserRoleUsernameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func (suite *UserRoleTestSuite) TestValidate_WithEmptyRole_ReturnsUserRoleEmptyRole() {
	//arrange
	suite.UserRole.Role = ""

	//act
	verr := suite.UserRole.Validate()

	//assert
	suite.Equal(models.ValidateUserRoleEmptyRole, verr)
}

func (suite *UserRoleTestSuite) TestValidate_RoleMaxLengthTestCases() {
	var role string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.UserRole.Role = role

		//act
		verr := suite.UserRole.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	role = helpers.CreateStringOfLength(models.UserRoleRoleMaxLength)
	expectedValidateError = models.ValidateUserRoleValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	role += "a"
	expectedValidateError = models.ValidateUserRoleRoleTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func TestUserRoleTestSuite(t *testing.T) {
	suite.Run(t, &UserRoleTestSuite{})
}