package models_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	Client *models.Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.Client = models.CreateNewClient("name", "redirect.com")
}

func (suite *ClientTestSuite) TestCreateNewClient_CreatesClientWithSuppliedFields() {
	//act
	name := "Client Name"
	url := "Redirect URL"
	client := models.CreateNewClient(name, url)

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

	name = helpers.CreateStringOfLength(30)
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

	url = helpers.CreateStringOfLength(100)
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

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
