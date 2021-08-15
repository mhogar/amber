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

func (suite *ClientCRUDTestSuite) TestCreateClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := models.CreateNewClient("")

	//act
	err := suite.Tx.CreateClient(client)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := models.CreateNewClient("")

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
	suite.CreateClient(client)

	client.Name = "new name"

	//act
	res, err := suite.Tx.UpdateClient(client)
	suite.Require().NoError(err)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultClient, err := suite.Tx.GetClientByUID(client.UID)
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_WhereClientIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.DeleteClient(uuid.New())

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_DeletesClientWithId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CreateClient(client)

	//act
	res, err := suite.Tx.DeleteClient(client.UID)
	suite.Require().NoError(err)

	//assert
	resultClient, err := suite.Tx.GetClientByUID(client.UID)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultClient)
}

func (suite *ClientCRUDTestSuite) TestGetClientByUId_WhereClientNotFound_ReturnsNilClient() {
	//act
	client, err := suite.Tx.GetClientByUID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(client)
}

func (suite *ClientCRUDTestSuite) TestGetClientByUId_GetsTheClientWithUId() {
	//arrange
	client := models.CreateNewClient("name")
	suite.CreateClient(client)

	//act
	resultClient, err := suite.Tx.GetClientByUID(client.UID)

	//assert
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func TestClientCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ClientCRUDTestSuite{})
}
