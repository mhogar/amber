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
	_, err := suite.Tx.UpdateClient(client)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_WhereClientIsNotFound_ReturnsFalseResult() {
	//arrange
	client := models.CreateNewClient("name")

	//act
	res, err := suite.Tx.UpdateClient(client)

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_UpdatesClientWithId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.SaveClient(client)

	client.Name = "new name"

	//act
	res, err := suite.Tx.UpdateClient(client)
	suite.Require().NoError(err)

	//assert
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	suite.True(res)
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_WhereClientIsNotFound_ReturnsFalseResult() {
	//arrange
	id := uuid.New()

	//act
	res, err := suite.Tx.DeleteClient(id)

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_DeletesClientWithId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.SaveClient(client)

	//act
	res, err := suite.Tx.DeleteClient(client.ID)
	suite.Require().NoError(err)

	//assert
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultClient)
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
