package runner_test

import (
	"authserver/testing/helpers"
	"authserver/tools/migration_runner/runner"
	"authserver/tools/migration_runner/runner/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	helpers.CustomSuite
	helpers.ScopeFactorySuite
	MigrationRunnerFactoryMock mocks.MigrationRunnerFactory
	MigrationRunnerMock        mocks.MigrationRunner
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()

	suite.MigrationRunnerFactoryMock = mocks.MigrationRunnerFactory{}
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
	err := runner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, false)

	//assert
	suite.NoError(err)

	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.MigrationRunnerFactoryMock.AssertCalled(suite.T(), "CreateMigrationRunner", &suite.DataExecutorMock)
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")
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
	err := runner.Run(&suite.ScopeFactoryMock, &suite.MigrationRunnerFactoryMock, true)

	//assert
	suite.NoError(err)

	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.MigrationRunnerFactoryMock.AssertCalled(suite.T(), "CreateMigrationRunner", &suite.DataExecutorMock)
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}
