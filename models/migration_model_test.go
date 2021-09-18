package models_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationTestSuite struct {
	helpers.CustomSuite
	Migration *models.Migration
}

func (suite *MigrationTestSuite) SetupTest() {
	suite.Migration = models.CreateMigration(
		"001",
	)
}

func (suite *MigrationTestSuite) TestCreateNewMigration_CreatesMigrationWithSuppliedFields() {
	//arrange
	timestamp := "this is a timestamp"

	//act
	migration := models.CreateMigration(timestamp)

	//assert
	suite.Require().NotNil(migration)
	suite.EqualValues(timestamp, migration.Timestamp)
}

func (suite *MigrationTestSuite) TestValidate_WithValidMigration_ReturnsModelValid() {
	//act
	verr := suite.Migration.Validate()

	//assert
	suite.EqualValues(models.ValidateMigrationValid, verr)
}

func (suite *MigrationTestSuite) TestValidate_WithVariousInvalidTimestamps_ReturnsError() {
	var timestamp string
	testCase := func() {
		//arrange
		suite.Migration.Timestamp = timestamp

		//act
		verr := suite.Migration.Validate()

		//assert
		suite.EqualValues(models.ValidateMigrationInvalidTimestamp, verr)
	}

	timestamp = "00"
	suite.Run("TooFewDigits", testCase)

	timestamp = "0000"
	suite.Run("TooManyDigits", testCase)

	timestamp = "0a0"
	suite.Run("ContainsNonDigit", testCase)
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, &MigrationTestSuite{})
}
