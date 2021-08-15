package controllers_test

import (
	"authserver/controllers"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	passwordhelpermocks "authserver/controllers/password_helpers/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenControllerTestSuite struct {
	ControllerTestSuite
	PasswordHasherMock passwordhelpermocks.PasswordHasher
	TokenController    controllers.CoreTokenController
}

func (suite *TokenControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.PasswordHasherMock = passwordhelpermocks.PasswordHasher{}
	suite.TokenController = controllers.CoreTokenController{
		PasswordHasher: &suite.PasswordHasherMock,
	}
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WithErrorGettingClientByID_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WhereClientWithIDisNotFound_ReturnsInvalidClient() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_client", "")
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WhereUserWithUsernameIsNotFound_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_grant", "username", "password")
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WherePasswordDoesNotMatch_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_grant", "username", "password")
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WithErrorSavingAccessToken_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Return(errors.New(""))

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenFromPassword_WithNoErrors_ReturnsNoError() {
	//arrange
	username := "username"
	password := "password"
	clientUID := uuid.New()

	client := models.CreateClient(clientUID, "name")
	user := models.CreateNewUser(username, []byte(password))

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(client, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Return(nil)

	//act
	token, rerr := suite.TokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientUID)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetClientByID", clientUID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", username)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", mock.Anything, password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveAccessToken", token)

	suite.Require().NotNil(token)
	suite.Equal(client, token.Client)
	suite.Equal(user, token.User)

	helpers.AssertOAuthNoError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithErrorDeletingAccessToken_ReturnsInternalError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.TokenController.DeleteToken(&suite.CRUDMock, token)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestDeleteToken_WithNoErrors_ReturnsNoError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(nil)

	//act
	rerr := suite.TokenController.DeleteToken(&suite.CRUDMock, token)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAccessToken", token)

	helpers.AssertNoError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestDeleteAllOtherUserTokens_WithErrorDeletingTokens_ReturnsInternalError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAllOtherUserTokens", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.TokenController.DeleteAllOtherUserTokens(&suite.CRUDMock, token)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *TokenControllerTestSuite) TestDeleteAllOtherUserTokens_WithNoErrors_ReturnsNoError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAllOtherUserTokens", mock.Anything).Return(nil)

	//act
	rerr := suite.TokenController.DeleteAllOtherUserTokens(&suite.CRUDMock, token)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllOtherUserTokens", token)

	helpers.AssertNoError(&suite.Suite, rerr)
}

func TestTokenControlTestSuite(t *testing.T) {
	suite.Run(t, &TokenControllerTestSuite{})
}
