package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/controllers/mocks"
	passwordhelpers "authserver/controllers/password_helpers"
	passwordhelpermocks "authserver/controllers/password_helpers/mocks"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	ControllerTestSuite
	PasswordHasherMock            passwordhelpermocks.PasswordHasher
	PasswordCriteriaValidatorMock passwordhelpermocks.PasswordCriteriaValidator
	ControllersMock               mocks.Controllers
	UserController                controllers.CoreUserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.PasswordHasherMock = passwordhelpermocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = passwordhelpermocks.PasswordCriteriaValidator{}
	suite.ControllersMock = mocks.Controllers{}
	suite.UserController = controllers.CoreUserController{
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
		AuthController:            &suite.ControllersMock,
	}
}

func (suite *UserControllerTestSuite) runValidateUserTestCases(validateFunc func(user *models.User) common.CustomError) {
	suite.Run("EmptyUsername_ReturnsClientError", func() {
		//arrange
		user := models.CreateUser("", []byte("password"), 0)

		//act
		cerr := validateFunc(user)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be empty")
	})

	suite.Run("UsernameTooLong_ReturnsClientError", func() {
		//arrange
		user := models.CreateUser(helpers.CreateStringOfLength(models.UserUsernameMaxLength+1), []byte("password"), 0)

		//act
		cerr := validateFunc(user)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be longer", fmt.Sprint(models.UserUsernameMaxLength))
	})

	suite.Run("InvalidRank_ReturnsClientError", func() {
		//arrange
		user := models.CreateUser("username", []byte("password"), -1)

		//act
		cerr := validateFunc(user)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "rank", "invalid")
	})
}

func (suite *UserControllerTestSuite) TestCreateUser_ValidateUserTestCases() {
	suite.runValidateUserTestCases(func(user *models.User) common.CustomError {
		resUser, cerr := suite.UserController.CreateUser(&suite.CRUDMock, user.Username, "password", user.Rank)
		suite.Nil(resUser)

		return cerr
	})
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, "username", "password", 0)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithNonUniqueUsername_ReturnsClientError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, "username", "password", 0)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "error creating user")
}

func (suite *UserControllerTestSuite) TestCreateUser_WherePasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, "username", "password", 0)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, "username", "password", 0)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorCreatingUser_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("CreateUser", mock.Anything).Return(errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, "username", "password", 0)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"
	password := "password"
	rank := 0

	hash := []byte("password hash")

	suite.CRUDMock.On("GetUserByUsername", username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.CRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password, rank)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", password)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateUser", user)

	suite.Require().NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(hash, user.PasswordHash)
	suite.Equal(rank, user.Rank)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalError() {
	//arrange
	user := models.CreateUser("username", []byte("password hash"), 0)
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserController.DeleteUser(&suite.CRUDMock, user.Username)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithFalseResultDeletingUser_ReturnsClientError() {
	//arrange
	username := "username"
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(false, nil)

	//act
	cerr := suite.UserController.DeleteUser(&suite.CRUDMock, username)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "user with username", username, "not found")
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithNoErrors_ReturnsNoError() {
	//arrange
	user := models.CreateUser("username", []byte("password hash"), 0)
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.DeleteUser(&suite.CRUDMock, user.Username)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUser", user.Username)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithClientErrorAuthenticatingUser_ReturnsClientError() {
	//arrange
	username := "username"
	oldPassword := "old password"
	newPassword := "new password"

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, username, oldPassword, newPassword)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "old password", "invalid")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithNonClientErrorAuthenticatingUser_ReturnsError() {
	//arrange
	username := "username"
	oldPassword := "old password"
	newPassword := "new password"

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, username, oldPassword, newPassword)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	user := models.CreateUser("username", nil, 0)

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user.Username, oldPassword, newPassword)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	user := models.CreateUser("username", nil, 0)

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user.Username, oldPassword, newPassword)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorUpdatingUser_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	user := models.CreateUser("username", nil, 0)

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user.Username, oldPassword, newPassword)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithNoErrors_ReturnsNoError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")

	user := models.CreateUser("username", oldPasswordHash, 0)

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(user, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user.Username, oldPassword, newPassword)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, user.Username, oldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", newPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", newPassword)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUser", user)

	suite.Equal(newPasswordHash, user.PasswordHash)
	helpers.AssertNoError(&suite.Suite, cerr)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
