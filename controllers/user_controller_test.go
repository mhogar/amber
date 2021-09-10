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
		user := models.CreateUser("", 0, nil)

		//act
		cerr := validateFunc(user)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be empty")
	})

	suite.Run("UsernameTooLong_ReturnsClientError", func() {
		//arrange
		user := models.CreateUser(helpers.CreateStringOfLength(models.UserUsernameMaxLength+1), 0, nil)

		//act
		cerr := validateFunc(user)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "username", "cannot be longer", fmt.Sprint(models.UserUsernameMaxLength))
	})

	suite.Run("InvalidRank_ReturnsClientError", func() {
		//arrange
		user := models.CreateUser("username", -1, nil)

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
	helpers.AssertClientError(&suite.Suite, cerr, "username", "already in use")
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

func (suite *UserControllerTestSuite) TestUpdateUser_ValidateUserTestCases() {
	suite.runValidateUserTestCases(func(user *models.User) common.CustomError {
		resUser, cerr := suite.UserController.UpdateUser(&suite.CRUDMock, user.Username, user.Rank)
		suite.Nil(resUser)

		return cerr
	})
}

func (suite *UserControllerTestSuite) TestUpdateUser_WithErrorUpdatingUser_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("UpdateUser", mock.Anything, mock.Anything).Return(false, errors.New(""))

	//act
	user, cerr := suite.UserController.UpdateUser(&suite.CRUDMock, "username", 0)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUser_WithFalseResultUpdatingUser_ReturnsClientError() {
	//arrange
	username := "username"
	suite.CRUDMock.On("UpdateUser", mock.Anything, mock.Anything).Return(false, nil)

	//act
	user, cerr := suite.UserController.UpdateUser(&suite.CRUDMock, username, 0)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "user with username", username, "not found")
}

func (suite *UserControllerTestSuite) TestUpdateUser_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"
	rank := 0

	suite.CRUDMock.On("UpdateUser", mock.Anything, mock.Anything).Return(true, nil)

	//act
	user, cerr := suite.UserController.UpdateUser(&suite.CRUDMock, username, rank)

	//assert
	suite.Require().NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(rank, user.Rank)

	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUser", user)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithClientErrorAuthenticatingUser_ReturnsClientError() {
	//arrange
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, "username", "old password", "new password")

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "old password", "incorrect")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithNonClientErrorAuthenticatingUser_ReturnsError() {
	//arrange
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, "username", "old password", "new password")

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, "username", "old password", "new password")

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "password", "not", "minimum criteria")
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, "username", "old password", "new password")

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithErrorUpdatingUserPassword_ReturnsInternalError() {
	//arrange
	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("UpdateUserPassword", mock.Anything, mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, "username", "old password", "new password")

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestUpdateUserPassword_WithNoErrors_ReturnsNoError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	username := "username"
	newPasswordHash := []byte("hashed new password")

	suite.ControllersMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.CRUDMock.On("UpdateUserPassword", mock.Anything, mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.UpdateUserPassword(&suite.CRUDMock, username, oldPassword, newPassword)

	//assert
	suite.ControllersMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, username, oldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", newPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", newPassword)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUserPassword", username, newPasswordHash)

	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalError() {
	//arrange
	user := models.CreateUser("username", 0, nil)
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
	user := models.CreateUser("username", 0, nil)
	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(true, nil)

	//act
	cerr := suite.UserController.DeleteUser(&suite.CRUDMock, user.Username)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUser", user.Username)
	helpers.AssertNoError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestVerifyUserRank_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	_, cerr := suite.UserController.VerifyUserRank(&suite.CRUDMock, "username", 0)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *UserControllerTestSuite) TestVerifyUserRank_WhereUserNotFounc_ReturnsClientError() {
	//arrange
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	_, cerr := suite.UserController.VerifyUserRank(&suite.CRUDMock, "username", 0)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "user", "not found")
}

func (suite *UserControllerTestSuite) TestVerifyUserRank_WithNoErrors_TestCases() {
	//arrange
	user := models.CreateUser("username", 5, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(user, nil)

	var rank int
	expectedResult := false

	testCase := func() {
		//act
		res, cerr := suite.UserController.VerifyUserRank(&suite.CRUDMock, user.Username, rank)

		//assert
		helpers.AssertNoError(&suite.Suite, cerr)
		suite.Equal(expectedResult, res)
	}

	rank = 4
	suite.Run("RankLessThanUser_ReturnsFalseResult", testCase)

	rank = 5
	suite.Run("RankEqualToUser_ReturnsFalseResult", testCase)

	rank = 6
	expectedResult = true
	suite.Run("RankGreaterThanUser_ReturnsTrueResult", testCase)

	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", user.Username, mock.Anything)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, &UserControllerTestSuite{})
}
