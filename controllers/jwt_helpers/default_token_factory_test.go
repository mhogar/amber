package jwthelpers_test

import (
	"errors"
	"testing"

	"github.com/mhogar/amber/config"
	jwthelpers "github.com/mhogar/amber/controllers/jwt_helpers"
	"github.com/mhogar/amber/controllers/jwt_helpers/mocks"
	loadermocks "github.com/mhogar/amber/loaders/mocks"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DefaultTokenFactoryTestSuite struct {
	helpers.CustomSuite
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
	message := "load private key error"
	suite.DataLoaderMock.On("Load", mock.Anything).Return(nil, errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken("key.json", uuid.New(), "username", "role")

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *DefaultTokenFactoryTestSuite) TestCreateToken_WithErrorSigningToken_ReturnsError() {
	//arrange
	viper.Set("token", config.TokenConfig{})

	suite.DataLoaderMock.On("Load", mock.Anything).Return(nil, nil)

	message := "sign token error"
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return("", errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken("key.json", uuid.New(), "username", "role")

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *DefaultTokenFactoryTestSuite) TestCreateToken_WithNoErrors_ReturnsToken() {
	//arrange
	cfg := config.TokenConfig{
		DefaultIssuer: "issuer",
		Lifetime:      60,
	}
	viper.Set("token", cfg)

	uri := "key.json"
	clientUID := uuid.New()
	username := "username"
	role := "role"

	privateKey := []byte("private key")
	token := "this_is_a_signed_token"

	suite.DataLoaderMock.On("Load", mock.Anything).Return(privateKey, nil)
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return(token, nil)

	//act
	resultToken, err := suite.TokenFactory.CreateToken(uri, clientUID, username, role)

	//assert
	suite.NoError(err)
	suite.Equal(token, resultToken)

	suite.DataLoaderMock.AssertCalled(suite.T(), "Load", uri)
	suite.TokenSignerMock.AssertCalled(suite.T(), "SignToken", mock.MatchedBy(func(tk *jwt.Token) bool {
		claims := tk.Claims.(jwthelpers.DefaultClaims)
		return claims.Username == username &&
			claims.Role == role &&
			claims.Audience == clientUID.String() &&
			claims.Issuer == cfg.DefaultIssuer &&
			claims.ExpiresAt-claims.IssuedAt == cfg.Lifetime
	}), privateKey)
}

func TestDefaultTokenFactoryTestSuite(t *testing.T) {
	suite.Run(t, &DefaultTokenFactoryTestSuite{})
}
