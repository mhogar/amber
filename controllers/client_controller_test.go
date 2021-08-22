package controllers_test

import (
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

func (suite *ClientControllerTestSuite) TestCreateClient_WithEmptyName_ReturnsClientError() {
	//arrange
	name := ""

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name, "")

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	name := helpers.CreateStringOfLength(31)

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name, "")

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithEmptyRedirectUrl_ReturnsClientError() {
	//arrange
	url := ""

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, "name", url)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithRedirectUrlGreaterThanMax_ReturnsClientError() {
	//arrange
	url := helpers.CreateStringOfLength(101)

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, "name", url)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be longer", fmt.Sprint(models.ClientRedirectUrlMaxLength))
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithErrorSavingClient_ReturnsInternalError() {
	//arrange
	name := "name"
	url := "redirect.com"

	suite.CRUDMock.On("CreateClient", mock.Anything).Return(errors.New(""))

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name, url)

	//assert
	suite.Nil(client)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	name := "name"
	url := "redirect.com"

	suite.CRUDMock.On("CreateClient", mock.Anything).Return(nil)

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name, url)

	//assert
	suite.Require().NotNil(client)
	suite.Equal(name, client.Name)
	suite.Equal(url, client.RedirectUrl)

	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateClient", client)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithEmptyName_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("", "")

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient(helpers.CreateStringOfLength(31), "")

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithEmptyRedirectUrl_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", "")

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithRedirectUrlGreaterThanMax_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", helpers.CreateStringOfLength(101))

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client redirect url", "cannot be longer", fmt.Sprint(models.ClientRedirectUrlMaxLength))
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithErrorUpdatingClient_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithFalseResultUpdatingClient_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, nil)

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client with id", client.UID.String(), "not found")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com")
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
