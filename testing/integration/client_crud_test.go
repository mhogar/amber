package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *ClientCRUDTestSuite) TestSaveClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := &models.Client{
		ID: uuid.Nil,
	}

	//act
	err := suite.Tx.SaveClient(client)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := &models.Client{
		ID: uuid.Nil,
	}

	//act
	err := suite.Tx.UpdateClient(client)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_UpdatesClientWithId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.SaveClient(client)

	client.Name = "new name"

	//act
	err := suite.Tx.UpdateClient(client)
	suite.Require().NoError(err)

	//assert
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func (suite *ClientCRUDTestSuite) TestGetClientById_WhereClientNotFound_ReturnsNilClient() {
	//act
	client, err := suite.Tx.GetClientByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(client)
}

func (suite *ClientCRUDTestSuite) TestGetClientById_GetsTheClientWithId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.SaveClient(client)

	//act
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func TestClientCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ClientCRUDTestSuite{})
}
