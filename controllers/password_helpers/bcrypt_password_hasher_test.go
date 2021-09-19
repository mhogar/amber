package passwordhelpers_test

import (
	"testing"

	passwordhelpers "github.com/mhogar/amber/controllers/password_helpers"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/stretchr/testify/suite"
)

type BCryptPasswordHasherTestSuite struct {
	helpers.CustomSuite
	BCryptPasswordHasher passwordhelpers.BCryptPasswordHasher
}

func (suite *BCryptPasswordHasherTestSuite) SetupTest() {
	suite.BCryptPasswordHasher = passwordhelpers.BCryptPasswordHasher{}
}

func (suite *BCryptPasswordHasherTestSuite) TestHashPassword_WithNoError_ReturnsHashAndNilError() {
	//act
	hash, err := suite.BCryptPasswordHasher.HashPassword("password")

	//assert
	suite.NotNil(hash)
	suite.NoError(err)
}

func (suite *BCryptPasswordHasherTestSuite) TestComparePasswords_WherePasswordDoesNotMatchHash_ReturnsError() {
	//act
	err := suite.BCryptPasswordHasher.ComparePasswords([]byte("incorrect hash"), "password")

	//assert
	suite.Error(err)
}

func (suite *BCryptPasswordHasherTestSuite) TestComparePasswords_WherePasswordMatchesHash_ReturnsNilError() {
	//arrange
	password := "password"
	hash, err := suite.BCryptPasswordHasher.HashPassword(password)
	suite.NoError(err)

	//act
	err = suite.BCryptPasswordHasher.ComparePasswords(hash, password)

	//assert
	suite.NoError(err)
}

func TestBCryptPasswordHasherTestSuite(t *testing.T) {
	suite.Run(t, &BCryptPasswordHasherTestSuite{})
}
