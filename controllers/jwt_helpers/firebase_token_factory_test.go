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

type FirebaseTokenFactoryTestSuite struct {
	helpers.CustomSuite
	JSONLoaderMock  loadermocks.JSONLoader
	TokenSignerMock mocks.TokenSigner
	TokenFactory    jwthelpers.FirebaseTokenFactory
}

func (suite *FirebaseTokenFactoryTestSuite) SetupTest() {
	suite.JSONLoaderMock = loadermocks.JSONLoader{}
	suite.TokenSignerMock = mocks.TokenSigner{}

	suite.TokenFactory = jwthelpers.FirebaseTokenFactory{
		JSONLoader:  &suite.JSONLoaderMock,
		TokenSigner: &suite.TokenSignerMock,
	}
}

func (suite *FirebaseTokenFactoryTestSuite) TestCreateToken_WithErrorLoadingJSON_ReturnsError() {
	//arrange
	message := "load service json error"
	suite.JSONLoaderMock.On("Load", mock.Anything, mock.Anything).Return(errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken("key.json", uuid.New(), "username", "role")

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *FirebaseTokenFactoryTestSuite) TestCreateToken_WithErrorSigningToken_ReturnsError() {
	//arrange
	viper.Set("token", config.TokenConfig{})

	suite.JSONLoaderMock.On("Load", mock.Anything, mock.Anything).Return(nil)

	message := "sign token error"
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return("", errors.New(message))

	//act
	token, err := suite.TokenFactory.CreateToken("key.json", uuid.New(), "username", "role")

	//assert
	suite.Empty(token)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *FirebaseTokenFactoryTestSuite) TestCreateToken_WithNoErrors_ReturnsToken() {
	//arrange
	cfg := config.TokenConfig{
		Lifetime: 60,
	}
	viper.Set("token", cfg)

	uri := "key.json"
	username := "username"
	role := "role"
	token := "this_is_a_signed_token"

	serviceJSON := jwthelpers.FirebaseServiceJSON{
		ClientEmail: "email",
		PrivateKey:  "private_key",
	}

	suite.JSONLoaderMock.On("Load", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		*args.Get(1).(*jwthelpers.FirebaseServiceJSON) = serviceJSON
	})
	suite.TokenSignerMock.On("SignToken", mock.Anything, mock.Anything).Return(token, nil)

	//act
	resultToken, err := suite.TokenFactory.CreateToken(uri, uuid.Nil, username, role)

	//assert
	suite.NoError(err)
	suite.Equal(token, resultToken)

	suite.JSONLoaderMock.AssertCalled(suite.T(), "Load", uri, mock.Anything)
	suite.TokenSignerMock.AssertCalled(suite.T(), "SignToken", mock.MatchedBy(func(tk *jwt.Token) bool {
		claims := tk.Claims.(jwthelpers.FirebaseClaims)
		return claims.UID == username &&
			claims.Issuer == serviceJSON.ClientEmail &&
			claims.Subject == serviceJSON.ClientEmail &&
			claims.ExpiresAt-claims.IssuedAt == cfg.Lifetime &&
			claims.Claims["role"] == role
	}), []byte(serviceJSON.PrivateKey))
}

func TestFirebaseTokenFactoryTestSuite(t *testing.T) {
	suite.Run(t, &FirebaseTokenFactoryTestSuite{})
}
