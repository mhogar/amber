package integration_test

import (
	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/dependencies"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/spf13/viper"
)

type CRUDTestSuite struct {
	helpers.CustomSuite
	Adapter  data.DataAdapter
	Executor data.DataExecutor
}

func (suite *CRUDTestSuite) SetupSuite() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	viper.Set("db_key", "integration")

	//-- create and setup the adapter --
	suite.Adapter = dependencies.ResolveDataAdapter()

	err = suite.Adapter.Setup()
	suite.Require().NoError(err)

	suite.Executor = suite.Adapter.GetExecutor()
}

func (suite *CRUDTestSuite) TearDownSuite() {
	suite.Adapter.CleanUp()
}

func (suite *CRUDTestSuite) SaveMigration(timestamp string) string {
	err := suite.Executor.CreateMigration(timestamp)
	suite.Require().NoError(err)

	return timestamp
}

func (suite *CRUDTestSuite) DeleteMigration(timestamp string) {
	err := suite.Executor.DeleteMigrationByTimestamp(timestamp)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveUser(user *models.User) *models.User {
	err := suite.Executor.CreateUser(user)
	suite.Require().NoError(err)

	return user
}

func (suite *CRUDTestSuite) DeleteUser(user *models.User) {
	_, err := suite.Executor.DeleteUser(user.Username)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveClient(client *models.Client) *models.Client {
	err := suite.Executor.CreateClient(client)
	suite.Require().NoError(err)

	return client
}

func (suite *CRUDTestSuite) DeleteClient(client *models.Client) {
	_, err := suite.Executor.DeleteClient(client.UID)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveUserRole(role *models.UserRole) *models.UserRole {
	err := suite.Executor.CreateUserRole(role)
	suite.Require().NoError(err)

	return role
}

func (suite *CRUDTestSuite) SaveSession(session *models.Session) *models.Session {
	err := suite.Executor.SaveSession(session)
	suite.Require().NoError(err)

	return session
}
