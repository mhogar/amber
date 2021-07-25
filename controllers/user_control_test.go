package controllers_test

import (
	"authserver/controllers"
	passwordhelpers "authserver/controllers/password_helpers"
	passwordhelpermocks "authserver/controllers/password_helpers/mocks"
	datamocks "authserver/data/mocks"
	"authserver/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControlTestSuite struct {
	suite.Suite
	CRUDMock                      datamocks.DataCRUD
	PasswordHasherMock            passwordhelpermocks.PasswordHasher
	PasswordCriteriaValidatorMock passwordhelpermocks.PasswordCriteriaValidator
	UserControl                   controllers.UserControl
}

func (suite *UserControlTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
	suite.PasswordHasherMock = passwordhelpermocks.PasswordHasher{}
	suite.PasswordCriteriaValidatorMock = passwordhelpermocks.PasswordCriteriaValidator{}
	suite.UserControl = controllers.UserControl{
		PasswordHasher:            &suite.PasswordHasherMock,
		PasswordCriteriaValidator: &suite.PasswordCriteriaValidatorMock,
	}
}

func (suite *UserControlTestSuite) TestCreateUser_WithEmptyUsername_ReturnsClientError() {
	//arrange
	username := ""
	password := "password"

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertClientError(&suite.Suite, rerr, "username cannot be empty")
}

func (suite *UserControlTestSuite) TestCreateUser_WithUsernameLongerThanMax_ReturnsClientError() {
	//arrange
	username := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //31 chars
	password := "password"

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertClientError(&suite.Suite, rerr, "username cannot be longer", fmt.Sprint(models.UserUsernameMaxLength))
}

func (suite *UserControlTestSuite) TestCreateUser_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestCreateUser_WithNonUniqueUsername_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertClientError(&suite.Suite, rerr, "error creating user")
}

func (suite *UserControlTestSuite) TestCreateUser_WherePasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertClientError(&suite.Suite, rerr, "password", "not", "minimum criteria")
}

func (suite *UserControlTestSuite) TestCreateUser_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestCreateUser_WithErrorCreatingUser_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("SaveUser", mock.Anything).Return(errors.New(""))

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestCreateUser_WithValidRequest_ReturnsOK() {
	//arrange
	username := "username"
	password := "password"

	hash := []byte("password hash")

	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(&models.AccessToken{}, nil)
	suite.CRUDMock.On("GetUserByUsername", username).Return(nil, nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(hash, nil)
	suite.CRUDMock.On("SaveUser", mock.Anything).Return(nil)

	//act
	user, rerr := suite.UserControl.CreateUser(&suite.CRUDMock, username, password)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", username)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", password)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveUser", user)

	suite.Require().NotNil(user)
	suite.Equal(username, user.Username)
	suite.Equal(hash, user.PasswordHash)

	AssertNoError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestDeleteUser_WithErrorDeletingUser_ReturnsInternalError() {
	//arrange
	user := models.CreateNewUser("username", []byte("password hash"))

	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.UserControl.DeleteUser(&suite.CRUDMock, user)

	//assert
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestDeleteUser_WithValidRequest_ReturnsOK() {
	//arrange
	user := models.CreateNewUser("username", []byte("password hash"))

	suite.CRUDMock.On("DeleteUser", mock.Anything).Return(nil)

	//act
	rerr := suite.UserControl.DeleteUser(&suite.CRUDMock, user)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteUser", user)

	AssertNoError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestUpdateUserPassword_WhereOldPasswordIsInvalid_ReturnsClientError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.UserControl.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	AssertClientError(&suite.Suite, rerr, "old password", "invalid")
}

func (suite *UserControlTestSuite) TestUpdateUserPassword_WhereNewPasswordDoesNotMeetCriteria_ReturnsClientError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaError(passwordhelpers.ValidatePasswordCriteriaTooShort, ""))

	//act
	rerr := suite.UserControl.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	AssertClientError(&suite.Suite, rerr, "password", "not", "minimum criteria")
}

func (suite *UserControlTestSuite) TestUpdateUserPassword_WithErrorHashingNewPassword_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, errors.New(""))

	//act
	rerr := suite.UserControl.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestUpdateUserPassword_WithErrorUpdatingUser_ReturnsInternalError() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"
	user := &models.User{}

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(nil, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.UserControl.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	AssertInternalError(&suite.Suite, rerr)
}

func (suite *UserControlTestSuite) TestUpdateUserPassword_WithValidRequest_ReturnsOK() {
	//arrange
	oldPassword := "old password"
	newPassword := "new password"

	oldPasswordHash := []byte("hashed old password")
	newPasswordHash := []byte("hashed new password")

	user := models.CreateNewUser("username", oldPasswordHash)

	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.PasswordCriteriaValidatorMock.On("ValidatePasswordCriteria", mock.Anything).Return(passwordhelpers.CreateValidatePasswordCriteriaValid())
	suite.PasswordHasherMock.On("HashPassword", mock.Anything).Return(newPasswordHash, nil)
	suite.CRUDMock.On("UpdateUser", mock.Anything).Return(nil)

	//act
	rerr := suite.UserControl.UpdateUserPassword(&suite.CRUDMock, user, oldPassword, newPassword)

	//assert
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", oldPasswordHash, oldPassword)
	suite.PasswordCriteriaValidatorMock.AssertCalled(suite.T(), "ValidatePasswordCriteria", newPassword)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "HashPassword", newPassword)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateUser", user)

	suite.Equal(newPasswordHash, user.PasswordHash)
	AssertNoError(&suite.Suite, rerr)
}

func TestUserControlTestSuite(t *testing.T) {
	suite.Run(t, &UserControlTestSuite{})
}
