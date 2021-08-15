package controllers_test

import (
	"authserver/controllers"
	"authserver/models"
	"authserver/testing/helpers"
	"errors"
	"fmt"
	"testing"

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
	client, rerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, rerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	name := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //31 chars

	//act
	client, rerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, rerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithErrorSavingClient_ReturnsInternalError() {
	//arrange
	name := "name"
	suite.CRUDMock.On("CreateClient", mock.Anything).Return(errors.New(""))

	//act
	client, rerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	name := "name"
	suite.CRUDMock.On("CreateClient", mock.Anything).Return(nil)

	//act
	client, rerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Require().NotNil(client)
	suite.Equal(name, client.Name)

	helpers.AssertNoError(&suite.Suite, rerr)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateClient", client)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithEmptyName_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("")

	//act
	rerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, rerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	//act
	rerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, rerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithErrorUpdatingClient_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, errors.New(""))

	//act
	rerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithFalseResultUpdatingClient_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, nil)

	//act
	rerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, rerr, "client with id", client.UID.String(), "not found")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(true, nil)

	//act
	rerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertNoError(&suite.Suite, rerr)
	suite.CRUDMock.AssertCalled(suite.T(), "UpdateClient", client)
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithErrorDeletingClient_ReturnsInternalError() {
	//arrange
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(false, errors.New(""))

	//act
	rerr := suite.ClientController.DeleteClient(&suite.CRUDMock, 0)

	//assert
	helpers.AssertInternalError(&suite.Suite, rerr)
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithFalseResultDeletingClient_ReturnsClientError() {
	//arrange
	id := int16(100)
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(false, nil)

	//act
	rerr := suite.ClientController.DeleteClient(&suite.CRUDMock, id)

	//assert
	helpers.AssertClientError(&suite.Suite, rerr, "client with id", fmt.Sprint(id), "not found")
}

func (suite *ClientControllerTestSuite) TestDeleteClient_WithNoErrors_ReturnsNoError() {
	//arrange
	id := int16(100)
	suite.CRUDMock.On("DeleteClient", mock.Anything).Return(true, nil)

	//act
	rerr := suite.ClientController.DeleteClient(&suite.CRUDMock, id)

	//assert
	helpers.AssertNoError(&suite.Suite, rerr)
	suite.CRUDMock.AssertCalled(suite.T(), "DeleteClient", id)
}

func TestClientControllerTestSuite(t *testing.T) {
	suite.Run(t, &ClientControllerTestSuite{})
}
