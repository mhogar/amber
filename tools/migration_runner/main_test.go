package main_test

import (
	datamocks "authserver/data/mocks"
	migrationrunner "authserver/tools/migration_runner"
	"authserver/tools/migration_runner/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	ScopeFactoryMock           datamocks.IScopeFactory
	MigrationRunnerFactoryMock mocks.IMigrationRunnerFactory
	MigrationRunnerMock        mocks.MigrationRunner
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.ScopeFactoryMock = datamocks.IScopeFactory{}
	suite.MigrationRunnerFactoryMock = mocks.IMigrationRunnerFactory{}
	suite.MigrationRunnerMock = mocks.MigrationRunner{}

	suite.MigrationRunnerFactoryMock.On("CreateMigrationRunner", mock.Anything).Return(&suite.MigrationRunnerMock)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownFalse_RunsUpMigration() {
	//arrange
	suite.MigrationRunnerMock.On("MigrateUp").Return(nil)

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, false)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")

	suite.NoError(err)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorRunningUpMigration_ReturnsError() {
	//arrange
	message := "MigrateUp test error"
	suite.MigrationRunnerMock.On("MigrateUp").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, false)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownTrue_RunsDownMigration() {
	//arrange
	suite.MigrationRunnerMock.On("MigrateDown").Return(nil)

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, true)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")

	suite.NoError(err)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorRunningDownMigration_ReturnsError() {
	//arrange
	message := "MigrateDown test error"
	suite.MigrationRunnerMock.On("MigrateDown").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, true)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}
