package main_test

import (
	"authserver/testing/helpers"
	migrationrunner "authserver/tools/migration_runner"
	"authserver/tools/migration_runner/interfaces/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	helpers.ScopeFactorySuite
	MigrationRunnerFactoryMock mocks.IMigrationRunnerFactory
	MigrationRunnerMock        mocks.MigrationRunner
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()

	suite.MigrationRunnerFactoryMock = mocks.IMigrationRunnerFactory{}
	suite.MigrationRunnerMock = mocks.MigrationRunner{}

	suite.MigrationRunnerFactoryMock.On("CreateMigrationRunner", mock.Anything).Return(&suite.MigrationRunnerMock)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownFalse_RunsUpMigrationAndReturnsResult() {
	//arrange
	message := "MigrateUp test error"
	suite.MigrationRunnerMock.On("MigrateUp").Return(errors.New(message))

	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(nil, func(err error) {
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, false)

	//assert
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.MigrationRunnerFactoryMock.AssertCalled(suite.T(), "CreateMigrationRunner", &suite.DataExecutorMock)
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")

	suite.NoError(err)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownTrue_RunsDownMigrationAndReturnsResult() {
	//arrange
	message := "MigrateDown test error"
	suite.MigrationRunnerMock.On("MigrateDown").Return(errors.New(message))

	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(nil, func(err error) {
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})

	//act
	err := migrationrunner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, true)

	//assert
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.MigrationRunnerFactoryMock.AssertCalled(suite.T(), "CreateMigrationRunner", &suite.DataExecutorMock)
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")

	suite.NoError(err)
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}