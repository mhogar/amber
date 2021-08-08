package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	Client *models.Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.Client = models.CreateNewClient("name")
}

func (suite *ClientTestSuite) TestCreateNewClient_CreatesClientWithSuppliedFields() {
	//act
	name := "Client Name"
	client := models.CreateNewClient(name)

	//assert
	suite.Require().NotNil(client)
	suite.NotEqual(uuid.Nil, client.ID)
	suite.Equal(name, client.Name)
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

	name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //30 chars
	expectedValidateError = models.ValidateClientValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	name += "a"
	expectedValidateError = models.ValidateClientNameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
