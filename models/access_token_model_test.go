package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccessTokenTestSuite struct {
	suite.Suite
	Token *models.AccessToken
}

func (suite *AccessTokenTestSuite) SetupTest() {
	suite.Token = models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient("name"),
	)
}

func (suite *AccessTokenTestSuite) TestCreateNewAccessToken_CreatesAccessTokenWithSuppliedFields() {
	//arrange
	user := models.CreateNewUser("", nil)
	client := models.CreateNewClient("name")

	//act
	token := models.CreateNewAccessToken(user, client)

	//assert
	suite.Require().NotNil(token)
	suite.NotEqual(token.ID, uuid.Nil)
	suite.Equal(token.User, user)
	suite.Equal(token.Client, client)
}

func (suite *AccessTokenTestSuite) TestValidate_WithValidAccessToken_ReturnsValid() {
	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenValid, verr)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilID_ReturnsAccessTokenInvalidID() {
	//arrange
	suite.Token.ID = uuid.Nil

	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilID, verr)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilUser_ReturnsAccessTokenNilUser() {
	//arrange
	suite.Token.User = nil

	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilUser, verr)
}

func (suite *AccessTokenTestSuite) TestValidate_WithInvalidUser_ReturnsAccessTokenInvalidUser() {
	//arrange
	suite.Token.User = models.CreateNewUser("", nil)

	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidUser, verr)
}

func (suite *AccessTokenTestSuite) TestValidate_WithNilClient_ReturnsAccessTokenNilClient() {
	//arrange
	suite.Token.Client = nil

	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenNilClient, verr)
}

func (suite *AccessTokenTestSuite) TestValidate_WithInvalidClient_ReturnsAccessTokenInvalidClient() {
	//arrange
	suite.Token.Client = &models.Client{ID: uuid.Nil}

	//act
	verr := suite.Token.Validate()

	//assert
	suite.Equal(models.ValidateAccessTokenInvalidClient, verr)
}

func TestAccessTokenTestSuite(t *testing.T) {
	suite.Run(t, &AccessTokenTestSuite{})
}
