package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	jwtmocks "authserver/controllers/jwt_helpers/mocks"
	"authserver/controllers/mocks"
	"authserver/models"
	"errors"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenControllerTestSuite struct {
	ControllerTestSuite
	ControllerMock           mocks.Controllers
	TokenFactorySelectorMock jwtmocks.TokenFactorySelector
	TokenFactoryMock         jwtmocks.TokenFactory
	TokenController          controllers.CoreTokenController
}

func (suite *TokenControllerTestSuite) SetupTest() {
	suite.ControllerTestSuite.SetupTest()

	suite.ControllerMock = mocks.Controllers{}
	suite.TokenFactorySelectorMock = jwtmocks.TokenFactorySelector{}
	suite.TokenFactoryMock = jwtmocks.TokenFactory{}

	suite.TokenController = controllers.CoreTokenController{
		AuthController:       &suite.ControllerMock,
		TokenFactorySelector: &suite.TokenFactorySelectorMock,
	}
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithErrorGettingClientByUID_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(nil, errors.New(""))

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, uuid.New(), "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WhereClientNotFound_ReturnsClientError() {
	//arrange
	clientUID := uuid.New()
	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(nil, nil)

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, clientUID, "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomClientError(cerr, "client with id", clientUID.String(), "not found")
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithClientErrorAuthenticatingUserWithPassword_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(""))

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomClientError(cerr, "invalid", "username", "password", "not assigned", "client")
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithNonClientErrorAuthenticatingUserWithPassword_ReturnsError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithErrorGettingUserRoleByUsernameAndClientUID_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WhereUserRoleForClientNotFound_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(nil, nil)

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, "username", "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomClientError(cerr, "invalid", "username", "password", "not assigned", "client")
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WhereTokenFactoryForTokenTypeNotFound_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	userRole := models.CreateUserRole(uuid.Nil, "username", "role")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(userRole, nil)
	suite.TokenFactorySelectorMock.On("Select", mock.Anything).Return(nil)

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, userRole.Username, "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithErrorCreatingToken_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	userRole := models.CreateUserRole(uuid.Nil, "username", "role")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(userRole, nil)
	suite.TokenFactorySelectorMock.On("Select", mock.Anything).Return(&suite.TokenFactoryMock)
	suite.TokenFactoryMock.On("CreateToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New(""))

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, userRole.Username, "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithErrorParsingRedirectUrl_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "invalid_\n_url", 0, "key.pem")
	userRole := models.CreateUserRole(uuid.Nil, "username", "role")

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(userRole, nil)
	suite.TokenFactorySelectorMock.On("Select", mock.Anything).Return(&suite.TokenFactoryMock)
	suite.TokenFactoryMock.On("CreateToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", nil)

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, userRole.Username, "password")

	//assert
	suite.Empty(tokenURL)
	suite.CustomInternalError(cerr)
}

func (suite *TokenControllerTestSuite) TestCreateTokenRedirectURL_WithNoErrors_ReturnsTokenRedirectURL() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	userRole := models.CreateUserRole(uuid.Nil, "username", "role")
	password := "password"
	token := "this_is_the_token_value"

	suite.CRUDMock.On("GetClientByUID", mock.Anything).Return(client, nil)
	suite.ControllerMock.On("AuthenticateUserWithPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil, common.NoError())
	suite.CRUDMock.On("GetUserRoleByClientUIDAndUsername", mock.Anything, mock.Anything).Return(userRole, nil)
	suite.TokenFactorySelectorMock.On("Select", mock.Anything).Return(&suite.TokenFactoryMock)
	suite.TokenFactoryMock.On("CreateToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(token, nil)

	//act
	tokenURL, cerr := suite.TokenController.CreateTokenRedirectURL(&suite.CRUDMock, client.UID, userRole.Username, password)

	//assert
	suite.CustomNoError(cerr)
	suite.Require().NotEmpty(tokenURL)

	url, err := url.Parse(tokenURL)
	suite.Require().NoError(err)
	suite.Equal(token, url.Query().Get("token"))

	suite.CRUDMock.AssertCalled(suite.T(), "GetClientByUID", client.UID)
	suite.ControllerMock.AssertCalled(suite.T(), "AuthenticateUserWithPassword", &suite.CRUDMock, userRole.Username, password)
	suite.CRUDMock.AssertCalled(suite.T(), "GetUserRoleByClientUIDAndUsername", client.UID, userRole.Username)
	suite.TokenFactorySelectorMock.AssertCalled(suite.T(), "Select", client.TokenType)
	suite.TokenFactoryMock.AssertCalled(suite.T(), "CreateToken", client.KeyUri, client.UID, userRole.Username, userRole.Role)
}

func TestTokenControllerTestSuite(t *testing.T) {
	suite.Run(t, &TokenControllerTestSuite{})
}
