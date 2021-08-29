package controllers_test

import (
	"authserver/common"
	"authserver/controllers"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ClientControllerTestSuite struct {
	ControllerTestSuite
	ClientController controllers.CoreClientController
}

func (suite *ClientControllerTestSuite) runValidateClientTestCases(validateFunc func(client *models.Client) common.CustomError) {
	suite.Run("EmptyName_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("", "", 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be empty")
	})

	suite.Run("NameGreaterThanMax_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient(helpers.CreateStringOfLength(models.ClientNameMaxLength+1), "", 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
	})

	suite.Run("EmptyRedirectUrl_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", "", 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be empty")
	})

	suite.Run("RedirectUrlGreaterThanMax_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", helpers.CreateStringOfLength(models.ClientRedirectUrlMaxLength+1), 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be longer", fmt.Sprint(models.ClientRedirectUrlMaxLength))
	})

	suite.Run("InvalidRedirectUrl_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", "invalid_\n_url", 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "invalid url")
	})

	suite.Run("InvalidTokenType_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", "redirect.com", -1, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client token type", "invalid")
	})

	suite.Run("EmptyKeyUri_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", "redirect.com", 0, "")

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client key uri", "cannot be empty")
	})

	suite.Run("KeyUriGreaterThanMax_ReturnsClientError", func() {
		//arrange
		client := models.CreateNewClient("name", "redirect.com", 0, helpers.CreateStringOfLength(models.ClientKeyUriMaxLength+1))

		//act
		cerr := validateFunc(client)

		//assert
		helpers.AssertClientError(&suite.Suite, cerr, "client key uri", "cannot be longer", fmt.Sprint(models.ClientKeyUriMaxLength))
	})
}

func (suite *ClientControllerTestSuite) TestCreateClient_ValidateClientTestCases() {
	suite.runValidateClientTestCases(func(client *models.Client) common.CustomError {
		resClient, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, client.Name, client.RedirectUrl, client.TokenType, client.KeyUri)
		suite.Nil(resClient)

		return cerr
	})
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithErrorSavingClient_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("CreateClient", mock.Anything).Return(errors.New(""))

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, "name", "redirect.com", 0, "key.pem")

	//assert
	suite.Nil(client)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	name := "name"
	url := "redirect.com"
	tokenType := 0
	uri := "key.pem"

	suite.CRUDMock.On("CreateClient", mock.Anything).Return(nil)

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name, url, tokenType, uri)

	//assert
	suite.Require().NotNil(client)
	suite.Equal(name, client.Name)
	suite.Equal(url, client.RedirectUrl)
	suite.Equal(tokenType, client.TokenType)
	suite.Equal(uri, client.KeyUri)

	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateClient", client)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_ValidateClientTestCases() {
	suite.runValidateClientTestCases(func(client *models.Client) common.CustomError {
		return suite.ClientController.UpdateClient(&suite.CRUDMock, client)
	})
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithErrorUpdatingClient_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithFalseResultUpdatingClient_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, nil)

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client with id", client.UID.String(), "not found")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(true, nil)

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateClient", client)
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithErrorDeletingClient_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.ClientController.DeleteClient(&suite.CRUDMock, uuid.Nil)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithFalseResultDeletingClient_ReturnsClientError() {
	//arrange
	uid := uuid.New()
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(false, nil)

	//act
	cerr := suite.ClientController.DeleteClient(&suite.CRUDMock, uid)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client with id", uid.String(), "not found")
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithNoErrors_ReturnsNoError() {
	//arrange
	uid := uuid.New()
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(true, nil)

	//act
	cerr := suite.ClientController.DeleteClient(&suite.CRUDMock, uid)

	//assert
	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteClient", uid)
}

func TestClientControllerTestSuite(t *testing.T) {
	suite.Run(t, &ClientControllerTestSuite{})
}
