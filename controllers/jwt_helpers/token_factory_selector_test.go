package jwthelpers_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TokenFactorySelectorTestSuite struct {
	suite.Suite
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

	tokenType = jwthelpers.TokenTypeFirebase
	expectedTokenFactory = &jwthelpers.FirebaseTokenFactory{}
	suite.Run("FirbaseTokenType_ReturnsFirebaseTokenFactory", testCase)
}

func TestTokenFactorySelectorTestSuite(t *testing.T) {
	suite.Run(t, &TokenFactorySelectorTestSuite{})
}
