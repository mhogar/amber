package jwthelpers_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"authserver/controllers/jwt_helpers/mocks"
	loadermocks "authserver/loaders/mocks"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DefaultTokenFactoryTestSuite struct {
	suite.Suite
	DataLoaderMock  loadermocks.RawDataLoader
	TokenSignerMock mocks.TokenSigner
	TokenFactory    jwthelpers.DefaultTokenFactory
}

func (suite *DefaultTokenFactoryTestSuite) SetupTest() {
	suite.DataLoaderMock = loadermocks.RawDataLoader{}
	suite.TokenSignerMock = mocks.TokenSigner{}

	suite.TokenFactory = jwthelpers.DefaultTokenFactory{
		DataLoader:  &suite.DataLoaderMock,
		TokenSigner: &suite.TokenSignerMock,
	}
}

func (suite *DefaultTokenFactoryTestSuite) TestCreateToken_WithErrorLoadingPrivateKey_ReturnsError() {
	//arrange
	uri := "key.json"
	username := "username"
	clientUID := uuid.New()

	message := "load private key error"
	suite.DataLoaderMock.On("Load", mock.Anything).Return(nil, errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken(uri, clientUID, username)

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *DefaultTokenFactoryTestSuite) TestCreateToken_WithErrorSigningToken_ReturnsError() {
	//arrange
	uri := "key.json"
	username := "username"
	clientUID := uuid.New()

	suite.DataLoaderMock.On("Load", mock.Anything).Return(nil, nil)

	message := "sign token error"
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return("", errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken(uri, clientUID, username)

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *DefaultTokenFactoryTestSuite) TestCreateToken_WithNoErrors_ReturnsToken() {
	//arrange
	uri := "key.json"
	username := "username"
	clientUID := uuid.New()

	privateKey := []byte("private key")
	token := "this_is_a_signed_token"

	suite.DataLoaderMock.On("Load", mock.Anything).Return(privateKey, nil)
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return(token, nil)

	//act
	resultToken, err := suite.TokenFactory.CreateToken(uri, clientUID, username)

	//assert
	suite.NoError(err)
	suite.Equal(token, resultToken)

	suite.DataLoaderMock.AssertCalled(suite.T(), "Load", uri)
	suite.TokenSignerMock.AssertCalled(suite.T(), "SignToken", mock.MatchedBy(func(tk *jwt.Token) bool {
		claims := tk.Claims.(jwthelpers.DefaultClaims)
		return claims.Username == username &&
			claims.Audience == clientUID.String() &&
			claims.ExpiresAt-claims.IssuedAt == 60
	}), privateKey)
}

func TestDefaultTokenFactoryTestSuite(t *testing.T) {
	suite.Run(t, &DefaultTokenFactoryTestSuite{})
}
