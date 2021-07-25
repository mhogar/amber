package data_integration_test

import (
	"authserver/common"
	"authserver/models"
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
	common.AssertError(&suite.Suite, err, "error", "client model")
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
	client := models.CreateNewClient()
	suite.SaveClient(suite.Tx, client)

	//act
	resultClient, err := suite.Tx.GetClientByID(client.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(client, resultClient)
}

func TestClientCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ClientCRUDTestSuite{})
}
