package integration_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *MigrationCRUDTestSuite) TestCreateMigration_WithInvalidTimestamp_ReturnsError() {
	//act
	err := suite.Executor.CreateMigration("invalid")

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "migration model")
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_WhereTimestampNotFound_ReturnsNilMigration() {
	//act
	migration, err := suite.Executor.GetMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
	suite.Nil(migration)
}

func (suite *MigrationCRUDTestSuite) TestGetMigrationByTimestamp_FindsMigration() {
	//arrange
	timestamp := suite.SaveMigration("999")

	//act
	migration, err := suite.Executor.GetMigrationByTimestamp(timestamp)

	//assert
	suite.NoError(err)
	suite.Require().NotNil(migration)
	suite.Equal(timestamp, migration.Timestamp)

	//clean up
	suite.DeleteMigration(timestamp)
}

func (suite *MigrationCRUDTestSuite) TestGetLatestTimestamp_ReturnsLatestTimestamp() {
	//arrange
	timestamps := []string{"991", "995", "992", "993"}
	for _, timestamp := range timestamps {
		suite.SaveMigration(timestamp)
	}

	//act
	timestamp, hasLatest, err := suite.Executor.GetLatestTimestamp()

	//assert
	suite.Equal(timestamps[1], timestamp)
	suite.True(hasLatest)
	suite.NoError(err)

	//clean up
	for _, timestamp := range timestamps {
		suite.DeleteMigration(timestamp)
	}
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_WithNoMigrationToDelete_ReturnsNilError() {
	//act
	err := suite.Executor.DeleteMigrationByTimestamp("DNE")

	//assert
	suite.NoError(err)
}

func (suite *MigrationCRUDTestSuite) TestDeleteMigrationByTimestamp_DeletesMigration() {
	//arrange
	timestamp := suite.SaveMigration("999")

	//act
	err := suite.Executor.DeleteMigrationByTimestamp(timestamp)

	//assert
	suite.Require().NoError(err)

	migration, err := suite.Executor.GetMigrationByTimestamp(timestamp)
	suite.NoError(err)
	suite.Nil(migration)
}

func TestMigrationCRUDTestSuite(t *testing.T) {
	suite.Run(t, &MigrationCRUDTestSuite{})
}
