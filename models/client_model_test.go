package models_test

import (
	"testing"

	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	helpers.CustomSuite
	Client *models.Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.Client = models.CreateNewClient("name", "redirect.com", 0, "key.pem")
}

func (suite *ClientTestSuite) TestCreateNewClient_CreatesClientWithSuppliedFields() {
	//arrange
	name := "Client Name"
	url := "Redirect URL"
	tokenType := 100
	uri := "Some Key URI"

	//act
	client := models.CreateNewClient(name, url, tokenType, uri)

	//assert
	suite.Require().NotNil(client)
	suite.NotEqual(uuid.Nil, client.UID)
	suite.Equal(name, client.Name)
	suite.Equal(url, client.RedirectUrl)
}

func (suite *ClientTestSuite) TestValidate_WithValidClient_ReturnsValid() {
	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientValid, verr)
}

func (suite *ClientTestSuite) TestValidate_WithNilUID_ReturnsClientNilUID() {
	//arrange
	suite.Client.UID = uuid.Nil

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientNilUID, verr)
}

func (suite *ClientTestSuite) TestValidate_WithEmptyName_ReturnsClientEmptyName() {
	//arrange
	suite.Client.Name = ""

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientEmptyName, verr)
}

func (suite *ClientTestSuite) TestValidate_ClientNameMaxLengthTestCases() {
	var name string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.Client.Name = name

		//act
		verr := suite.Client.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	name = helpers.CreateStringOfLength(models.ClientNameMaxLength)
	expectedValidateError = models.ValidateClientValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	name += "a"
	expectedValidateError = models.ValidateClientNameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func (suite *ClientTestSuite) TestValidate_WithEmptyRedirectUrl_ReturnsClientEmptyRedirectUrl() {
	//arrange
	suite.Client.RedirectUrl = ""

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientEmptyRedirectUrl, verr)
}

func (suite *ClientTestSuite) TestValidate_ClientRedirectUrlMaxLengthTestCases() {
	var url string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.Client.RedirectUrl = url

		//act
		verr := suite.Client.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	url = helpers.CreateStringOfLength(models.ClientRedirectUrlMaxLength)
	expectedValidateError = models.ValidateClientValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	url += "a"
	expectedValidateError = models.ValidateClientRedirectUrlTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func (suite *ClientTestSuite) TestValidate_WithInvalidRedirectUrl_ReturnsClientInvalidRedirectUrl() {
	//arrange
	suite.Client.RedirectUrl = "invalid_\n_url"

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientInvalidRedirectUrl, verr)
}

func (suite *ClientTestSuite) TestValidate_InvalidTokenTypeTestCases() {
	var tokenType int

	testCase := func() {
		//arrange
		suite.Client.TokenType = tokenType

		//act
		verr := suite.Client.Validate()

		//assert
		suite.Equal(models.ValidateClientInvalidTokenType, verr)
	}

	tokenType = -1
	suite.Run("LessThanSmallestTokenTypeValue", testCase)

	tokenType = 2
	suite.Run("MoreThanLargestTokenTypeValue", testCase)
}

func (suite *ClientTestSuite) TestValidate_WithEmptyKeyUri_ReturnsClientEmptyKeyUri() {
	//arrange
	suite.Client.KeyUri = ""

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientEmptyKeyUri, verr)
}

func (suite *ClientTestSuite) TestValidate_KeyUriMaxLengthTestCases() {
	var uri string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.Client.KeyUri = uri

		//act
		verr := suite.Client.Validate()

		//assert
		suite.Equal(expectedValidateError, verr)
	}

	uri = helpers.CreateStringOfLength(models.ClientKeyUriMaxLength)
	expectedValidateError = models.ValidateClientValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	uri += "a"
	expectedValidateError = models.ValidateClientKeyUriTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
