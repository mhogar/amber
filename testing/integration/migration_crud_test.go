package integration_test

import (
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//act
	err := suite.Tx.CreateMigration("invalid")

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "migration model")
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_WhereTimestampNotFound_ReturnsNilMigration() {
	//act
	migration, err := suite.Tx.GetMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(migration)
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_FindsMigration() {
	//arrange
	timestamp := "999"
	err := suite.Tx.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	migration, err := suite.Tx.GetMigrationByTimestamp(timestamp)

	//assert
	suite.NoError(err)
	suite.Require().NotNil(migration)
	suite.Equal(timestamp, migration.Timestamp)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_ReturnsLatestTimestamp() {
	//arrange
	timestamps := []string{
		"991",
		"995",
		"992",
		"993",
	}

	for _, timestamp := range timestamps {
		err := suite.Tx.CreateMigration(timestamp)
		suite.Require().NoError(err)
	}

	//act
	timestamp, hasLatest, err := suite.Tx.GetLatestTimestamp()

	//assert
	suite.Equal(timestamps[1], timestamp)
	suite.True(hasLatest)
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_WithNoMigrationToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_DeletesMigration() {
	//arrange
	timestamp := "999"
	err := suite.Tx.CreateMigration(timestamp)
	suite.Require().NoError(err)

	//act
	err = suite.Tx.DeleteMigrationByTimestamp(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Tx.GetMigrationByTimestamp(timestamp)
	suite.NoError(err)
	suite.Nil(migration)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
