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
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	name := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //31 chars

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithErrorSavingClient_ReturnsInternalError() {
	//arrange
	name := "name"
	suite.CRUDMock.On("CreateClient", mock.Anything).Return(errors.New(""))

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Nil(client)
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestCreateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	name := "name"
	suite.CRUDMock.On("CreateClient", mock.Anything).Return(nil)

	//act
	client, cerr := suite.ClientController.CreateClient(&suite.CRUDMock, name)

	//assert
	suite.Require().NotNil(client)
	suite.Equal(name, client.Name)

	helpers.AssertNoError(&suite.Suite, cerr)
	suite.CRUDMock.AssertCalled(suite.T(), "CreateClient", client)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithEmptyName_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("")

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be empty")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNameGreaterThanMax_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client name", "cannot be longer", fmt.Sprint(models.ClientNameMaxLength))
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithErrorUpdatingClient_ReturnsInternalError() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, errors.New(""))

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertInternalError(&suite.Suite, cerr)
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithFalseResultUpdatingClient_ReturnsClientError() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CRUDMock.On("UpdateClient", mock.Anything).Return(false, nil)

	//act
	cerr := suite.ClientController.UpdateClient(&suite.CRUDMock, client)

	//assert
	helpers.AssertClientError(&suite.Suite, cerr, "client with id", client.UID.String(), "not found")
}

func (suite *ClientControllerTestSuite) TestUpdateClient_WithNoErrors_ReturnsNoError() {
	//arrange
	client := models.CreateNewClient("name")
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
	helpers.AssertClientError(&suite.Suite, cerr, "client with uid", uid.String(), "not found")
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
