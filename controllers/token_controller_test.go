package controllers_test

import (
	"authserver/controllers"
	datamocks "authserver/data/mocks"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	passwordhelpermocks "authserver/controllers/password_helpers/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenControlTestSuite struct {
	suite.Suite
	CRUDMock            datamocks.DataCRUD
	PasswordHasherMock  passwordhelpermocks.PasswordHasher
	CoreTokenController controllers.CoreTokenController
}

func (suite *TokenControlTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
	suite.PasswordHasherMock = passwordhelpermocks.PasswordHasher{}
	suite.CoreTokenController = controllers.CoreTokenController{
		PasswordHasher: &suite.PasswordHasherMock,
	}
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WithErrorGettingClientByID_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WhereClientWithIDisNotFound_ReturnsInvalidClient() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_client", "")
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WithErrorGettingScopeByName_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WhereNoScopeWithNameisNotFound_ReturnsInvalidScope() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_scope", "")
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WithErrorGettingUserByUsername_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WhereUserWithUsernameIsNotFound_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_grant", "username", "password")
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WherePasswordDoesNotMatch_ReturnsClientError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(errors.New(""))

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthClientError(&suite.Suite, rerr, "invalid_grant", "username", "password")
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WithErrorSavingAccessToken_ReturnsInternalError() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scope := "scope"

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(&models.Client{}, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(&models.Scope{}, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(&models.User{}, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Return(errors.New(""))

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scope)

	//assert
	suite.Nil(token)
	helpers.AssertOAuthInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestCreateTokenFromPassword_WithValidRequest_ReturnsOK() {
	//arrange
	username := "username"
	password := "password"
	clientID := uuid.New()
	scopeName := "scope"

	client := &models.Client{ID: clientID}
	scope := &models.Scope{ID: uuid.New()}
	user := &models.User{ID: uuid.New()}

	suite.CRUDMock.On("GetClientByID", mock.Anything).Return(client, nil)
	suite.CRUDMock.On("GetScopeByName", mock.Anything).Return(scope, nil)
	suite.CRUDMock.On("GetUserByUsername", mock.Anything).Return(user, nil)
	suite.PasswordHasherMock.On("ComparePasswords", mock.Anything, mock.Anything).Return(nil)
	suite.CRUDMock.On("SaveAccessToken", mock.Anything).Return(nil)

	//act
	token, rerr := suite.CoreTokenController.CreateTokenFromPassword(&suite.CRUDMock, username, password, clientID, scopeName)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "GetClientByID", clientID)
	suite.CRUDMock.AssertCalled(suite.T(), "GetScopeByName", scopeName)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserByUsername", username)
	suite.PasswordHasherMock.AssertCalled(suite.T(), "ComparePasswords", mock.Anything, password)
	suite.CRUDMock.AssertCalled(suite.T(), "SaveAccessToken", token)

	suite.Require().NotNil(token)
	suite.Equal(client, token.Client)
	suite.Equal(scope, token.Scope)
	suite.Equal(user, token.User)

	helpers.AssertOAuthNoError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestDeleteToken_WithErrorDeletingAccessToken_ReturnsInternalError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.CoreTokenController.DeleteToken(&suite.CRUDMock, token)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestDeleteToken_WithValidRequest_ReturnsOK() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAccessToken", mock.Anything).Return(nil)

	//act
	rerr := suite.CoreTokenController.DeleteToken(&suite.CRUDMock, token)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAccessToken", token)

	helpers.AssertNoError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestDeleteAllOtherUserTokens_WithErrorDeletingTokens_ReturnsInternalError() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAllOtherUserTokens", mock.Anything).Return(errors.New(""))

	//act
	rerr := suite.CoreTokenController.DeleteAllOtherUserTokens(&suite.CRUDMock, token)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *TokenControlTestSuite) TestDeleteAllOtherUserTokens_WithValidRequest_ReturnsOK() {
	//arrange
	token := &models.AccessToken{}

	suite.CRUDMock.On("DeleteAllOtherUserTokens", mock.Anything).Return(nil)

	//act
	rerr := suite.CoreTokenController.DeleteAllOtherUserTokens(&suite.CRUDMock, token)

	//assert
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteAllOtherUserTokens", token)

	helpers.AssertNoError(&suite.Suite, rerr)
}

func TestTokenControlTestSuite(t *testing.T) {
	suite.Run(t, &TokenControlTestSuite{})
}
