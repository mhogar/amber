package integration_test

import (
	"testing"

	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *ClientCRUDTestSuite) TestCreateClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := models.CreateNewClient("", "", 0, "")

	//act
	err := suite.Executor.CreateClient(client)

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestGetClients_GetsClientsOrderedByName() {
	//arrange
	client1 := suite.SaveClient(models.CreateNewClient("name1", "redirect.com", 0, "key.pem"))
	client2 := suite.SaveClient(models.CreateNewClient("name2", "redirect.com", 0, "key.pem"))

	//act
	clients, err := suite.Executor.GetClients()

	//assert
	suite.NoError(err)

	suite.Require().Len(clients, 2)
	suite.EqualValues(clients[0], client1)
	suite.EqualValues(clients[1], client2)

	//clean up
	suite.DeleteClient(client1)
	suite.DeleteClient(client2)
}

func (suite *ClientCRUDTestSuite) TestGetClientByUID_WhereClientNotFound_ReturnsNilClient() {
	//act
	client, err := suite.Executor.GetClientByUID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(client)
}

func (suite *ClientCRUDTestSuite) TestGetClientByUId_GetsTheClientWithUID() {
	//arrange
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))

	//act
	resultClient, err := suite.Executor.GetClientByUID(client.UID)

	//assert
	suite.NoError(err)
	suite.EqualValues(client, resultClient)

	//clean up
	suite.DeleteClient(client)
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_WithInvalidClient_ReturnsError() {
	//arrange
	client := models.CreateNewClient("", "", 0, "")

	//act
	_, err := suite.Executor.UpdateClient(client)

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "client model")
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_WhereClientIsNotFound_ReturnsFalseResult() {
	//arrange
	client := models.CreateNewClient("name", "redirect.com", 0, "key.pem")

	//act
	res, err := suite.Executor.UpdateClient(client)

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *ClientCRUDTestSuite) TestUpdateClient_UpdatesClientWithId() {
	//arrange
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	client.Name = "new name"

	//act
	res, err := suite.Executor.UpdateClient(client)
	suite.Require().NoError(err)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	resultClient, err := suite.Executor.GetClientByUID(client.UID)
	suite.NoError(err)
	suite.EqualValues(client, resultClient)

	//clean up
	suite.DeleteClient(client)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_WhereClientIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Executor.DeleteClient(uuid.New())

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_DeletesClientWithId() {
	//arrange
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))

	//act
	res, err := suite.Executor.DeleteClient(client.UID)
	suite.Require().NoError(err)

	//assert
	resultClient, err := suite.Executor.GetClientByUID(client.UID)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultClient)
}

func (suite *ClientCRUDTestSuite) TestDeleteClient_AlsoDeletesAllRolesForClient() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	client := suite.SaveClient(models.CreateNewClient("name", "redirect.com", 0, "key.pem"))
	suite.SaveUserRole(models.CreateUserRole(client.UID, user.Username, "role"))

	//act
	res, err := suite.Executor.DeleteClient(client.UID)

	//assert
	suite.True(res)
	suite.Require().NoError(err)

	role, err := suite.Executor.GetUserRoleByClientUIDAndUsername(client.UID, user.Username)
	suite.NoError(err)
	suite.Nil(role)

	//clean up
	suite.DeleteUser(user)
}

func TestClientCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ClientCRUDTestSuite{})
}
