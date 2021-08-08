package integration_test

import (
	"authserver/config"
	"authserver/data"
	"authserver/dependencies"
	"authserver/models"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CRUDTestSuite struct {
	suite.Suite
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

func (suite *CRUDTestSuite) SaveUser(user *models.User) {
	err := suite.Tx.SaveUser(user)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveScope(scope *models.Scope) {
	err := suite.Tx.SaveScope(scope)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveClient(client *models.Client) {
	err := suite.Tx.SaveClient(client)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveAccessToken(token *models.AccessToken) {
	err := suite.Tx.SaveAccessToken(token)
	suite.Require().NoError(err)
}

func (suite *CRUDTestSuite) SaveAccessTokenAndFields(token *models.AccessToken) {
	suite.SaveUser(token.User)
	suite.SaveClient(token.Client)
	suite.SaveScope(token.Scope)
	suite.SaveAccessToken(token)
}
