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
	Adapter data.DataAdapter
	Tx      data.Transaction
}

func (suite *CRUDTestSuite) SetupSuite() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	viper.Set("db_key", "integration")

	//-- create and setup the adapter --
	suite.Adapter = dependencies.ResolveDataAdapter()

	err = suite.Adapter.Setup()
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) TearDownSuite() {
	suite.Adapter.CleanUp()
}

func (suite *CRUDTestSuite) SetupTest() {
	//start a new transaction for every test
	tx, err := suite.Adapter.GetExecutor().CreateTransaction()
	suite.Require().NoError(err)

	suite.Tx = tx
}

func (suite *CRUDTestSuite) TearDownTest() {
	//rollback the transaction after each test
	err := suite.Tx.Rollback()
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveUser(user *models.User) *models.User {
	err := suite.Tx.CreateUser(user)
	suite.Require().NoError(err)

	return user
}

func (suite *CRUDTestSuite) SaveClient(client *models.Client) *models.Client {
	err := suite.Tx.CreateClient(client)
	suite.Require().NoError(err)

	return client
}

func (suite *CRUDTestSuite) SaveUserRole(role *models.UserRole) *models.UserRole {
	err := suite.Tx.CreateUserRole(role)
	suite.Require().NoError(err)

	return role
}

func (suite *CRUDTestSuite) SaveSession(session *models.Session) *models.Session {
	err := suite.Tx.SaveSession(session)
	suite.Require().NoError(err)

	return session
}
