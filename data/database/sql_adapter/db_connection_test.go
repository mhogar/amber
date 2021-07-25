package sqladapter_test

import (
	"authserver/common"
	"authserver/config"
	sqladapter "authserver/data/database/sql_adapter"
	"authserver/dependencies"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DbConnectionTestSuite struct {
	suite.Suite
	Adapter *sqladapter.SQLAdapter
}

func (suite *DbConnectionTestSuite) SetupTest() {
	err := config.InitConfig("../..")
	suite.Require().NoError(err)

	suite.Adapter = sqladapter.CreateSQLAdpater("integration", dependencies.ResolveSQLDriver())
}

func (suite *DbConnectionTestSuite) TestOpenConnection_WhereConnectionStringIsNotFound_ReturnsError() {
	//arrange
	dbKey := "not a real dbkey"
	suite.Adapter.DbKey = dbKey

	//act
	err := suite.Adapter.OpenConnection()

	//assert
	common.AssertError(&suite.Suite, err, "no connection string", dbKey)
}

func (suite *DbConnectionTestSuite) TestCloseConnection_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.Adapter.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.Adapter.CloseConnection()

	//assert
	suite.NoError(err)
	suite.Nil(suite.Adapter.DB)
}

func (suite *DbConnectionTestSuite) TestPing_WithValidConnection_ReturnsNoError() {
	//arrange
	err := suite.Adapter.OpenConnection()
	suite.Require().NoError(err)

	//act
	err = suite.Adapter.Ping()

	//assert
	suite.NoError(err)
}

func TestDbConnectionTestSuite(t *testing.T) {
	suite.Run(t, &DbConnectionTestSuite{})
}
