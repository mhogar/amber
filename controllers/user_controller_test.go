package controllers_test

import (
	"authserver/controllers"
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
	UserController                controllers.CoreUserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.PasswordHasherMock = passwordhelpermocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = passwordhelpermocks.PasswordCriteriaValidator{}
	suite.UserController = controllers.CoreUserController{
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
	}
}

func (suite *UserControllerTestSuite) TestCreateUser_WithEmptyUsername_ReturnsClientError() {
	//arrange
	username := ""
	password := "password"

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "username cannot be empty")
}

func (suite *UserControllerTestSuite) TestCreateUser_WithUsernameLongerThanMax_ReturnsClientError() {
	//arrange
	username := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //31 chars
	password := "password"

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "username cannot be longer", fmt.Sprint(models.UserUsernameMaxLength))
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithNonUniqueUsername_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "error creating user")
}

func (suite *UserControllerTestSuite) TestCreateUser_WherePasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithErrorCreatingUser_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("CreateUser", mock.Anything).Return(errors.New(""))

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestCreateUser_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"
	password := "password"

	hash := []byte("password hash")

	suite.CRUDMock.On("GetUserByUsername", username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.CRUDMock.On("CreateUser", mock.Anything).Return(nil)

	//act
	user, cerr := suite.UserController.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", password)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateUser", user)

	suite.Require().NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(hash, user.PasswordHash)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalError() {
	//arrange
	user := models.CreateNewUser("username", []byte("password hash"))
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
	user := models.CreateNewUser("username", []byte("password hash"))
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.DeleteUser(&suite.CRUDMock, user.Username)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUser", user.Username)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WhereOldPasswordIsInvalid_ReturnsClientError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "old password", "invalid")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorUpdatingUser_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithNoErrors_ReturnsNoError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")

	user := models.CreateNewUser("username", oldPasswordHash)

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", oldPasswordHash, oldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", newPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", newPassword)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUser", user)

	suite.Equal(newPasswordHash, user.PasswordHash)
	helpers.AssertNoError(&suite.Suite, cerr)
}

func TestUserControlTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
