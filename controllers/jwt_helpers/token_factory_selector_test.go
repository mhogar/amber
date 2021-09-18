package jwthelpers_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TokenFactorySelectorTestSuite struct {
	helpers.CustomSuite
	TokenFactorySelector jwthelpers.CoreTokenFactorySelector
}

func (suite *TokenFactorySelectorTestSuite) SetupTest() {
	suite.TokenFactorySelector = jwthelpers.CoreTokenFactorySelector{}
}

func (suite *TokenFactorySelectorTestSuite) TestSelect_ChoosesCorrectTokeFactory_Tests() {
	var tokenType int
	var expectedTokenFactory jwthelpers.TokenFactory

	testCase := func() {
		//act
		tf := suite.TokenFactorySelector.Select(tokenType)

		//assert
		suite.IsType(expectedTokenFactory, tf)
	}

	tokenType = -1
	expectedTokenFactory = nil
	suite.Run("UnknownTokenType_ReturnsNil", testCase)

	tokenType = models.ClientTokenTypeDefault
	expectedTokenFactory = &jwthelpers.DefaultTokenFactory{}
	suite.Run("FirbaseTokenType_ReturnsDefaultTokenFactory", testCase)

	tokenType = models.ClientTokenTypeFirebase
	expectedTokenFactory = &jwthelpers.FirebaseTokenFactory{}
	suite.Run("FirbaseTokenType_ReturnsFirebaseTokenFactory", testCase)
}

func TestTokenFactorySelectorTestSuite(t *testing.T) {
	suite.Run(t, &TokenFactorySelectorTestSuite{})
}
