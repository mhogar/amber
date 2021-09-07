package controllers_test

import (
	"authserver/controllers"
	"authserver/controllers/password_helpers/mocks"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	ControllerTestSuite
	PasswordHasherMock mocks.PasswordHasher
	AuthController     controllers.CoreAuthController
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.PasswordHasherMock = mocks.PasswordHasher{}
	suite.AuthController = controllers.CoreAuthController{
		PasswordHasher: &suite.PasswordHasherMock,
	}
}

func (suite *AuthControllerTestSuite) TestAuthenticateUserWithPassword_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	user, cerr := suite.AuthController.AuthenticateUserWithPassword(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *AuthControllerTestSuite) TestAuthenticateUserWithPassword_WhereUserWithUsernameIsNotFound_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	user, cerr := suite.AuthController.AuthenticateUserWithPassword(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "invalid", "username", "password")
}

func (suite *AuthControllerTestSuite) TestAuthenticateUserWithPassword_WherePasswordDoesNotMatch_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	user, cerr := suite.AuthController.AuthenticateUserWithPassword(&suite.CRUDMock, username, password)

	//assert
	suite.Nil(user)
	helpers.AssertClientError(&suite.Suite, cerr, "invalid", "username", "password")
}

func (suite *AuthControllerTestSuite) TestAuthenticateUserWithPassword_WithNoErrors_ReturnsNoError() {
	//arrange
	password := "password"
	existingUser := models.CreateUser("username", 0, nil)

	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(existingUser, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)

	//act
	user, cerr := suite.AuthController.AuthenticateUserWithPassword(&suite.CRUDMock, existingUser.Username, password)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", existingUser.Username)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", existingUser.PasswordHash, password)

	suite.Equal(existingUser, user)
	helpers.AssertNoError(&suite.Suite, cerr)
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, &AuthControllerTestSuite{})
}
