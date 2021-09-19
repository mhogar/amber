package data_test

import (
	"errors"
	"testing"

	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/data/mocks"
	"github.com/mhogar/amber/testing/helpers"

	"github.com/stretchr/testify/suite"
)

type ScopeFactoryTestSuite struct {
	helpers.CustomSuite
	DataAdapterMock  mocks.DataAdapter
	DataExecutorMock mocks.DataExecutor
	TransactionMock  mocks.Transaction
	ScopeFactory     data.ScopeFactory
}

func (suite *ScopeFactoryTestSuite) SetupTest() {
	suite.DataAdapterMock = mocks.DataAdapter{}
	suite.DataExecutorMock = mocks.DataExecutor{}
	suite.TransactionMock = mocks.Transaction{}

	suite.ScopeFactory = data.CoreScopeFactory{
		DataAdapter: &suite.DataAdapterMock,
	}
}

func (suite *ScopeFactoryTestSuite) TestCreateDataExecutorScope_WithErrorSettingUpDataAdapter_ReturnsError() {
	//arrange
	message := "setup error"
	suite.DataAdapterMock.On("Setup").Return(errors.New(message))

	//act
	err := suite.ScopeFactory.CreateDataExecutorScope(func(_ data.DataExecutor) error {
		return nil
	})

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *ScopeFactoryTestSuite) TestCreateDataExecutorScope_ReturnsResultFromBody() {
	//arrange
	message := "body error"

	suite.DataAdapterMock.On("Setup").Return(nil)
	suite.DataAdapterMock.On("CleanUp").Return(nil)
	suite.DataAdapterMock.On("GetExecutor").Return(&suite.DataExecutorMock)

	//act
	err := suite.ScopeFactory.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		suite.Equal(&suite.DataExecutorMock, exec)
		return errors.New(message)
	})

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)

	suite.DataAdapterMock.AssertCalled(suite.T(), "Setup")
	suite.DataAdapterMock.AssertCalled(suite.T(), "GetExecutor")
	suite.DataAdapterMock.AssertCalled(suite.T(), "CleanUp")
}

func (suite *ScopeFactoryTestSuite) TestCreateTransactionScope_WithErrorCreatingTransaction_ReturnsError() {
	//arrange
	message := "create transaction error"
	suite.DataExecutorMock.On("CreateTransaction").Return(nil, errors.New(message))

	//act
	err := suite.ScopeFactory.CreateTransactionScope(&suite.DataExecutorMock, func(_ data.Transaction) (bool, error) {
		return true, nil
	})

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *ScopeFactoryTestSuite) TestCreateTransactionScope_WithErrorFromBody_ReturnsErrorAndRollsBackTransaction() {
	//arrange
	message := "body error"

	suite.DataExecutorMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.TransactionMock.On("Rollback").Return(nil)

	//act
	err := suite.ScopeFactory.CreateTransactionScope(&suite.DataExecutorMock, func(tx data.Transaction) (bool, error) {
		suite.Equal(&suite.TransactionMock, tx)
		return false, errors.New(message)
	})

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)

	suite.TransactionMock.AssertCalled(suite.T(), "Rollback")
}

func (suite *ScopeFactoryTestSuite) TestCreateTransactionScope_WithFailureFromBody_RollsBackTransaction() {
	//arrange
	suite.DataExecutorMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.TransactionMock.On("Rollback").Return(nil)

	//act
	err := suite.ScopeFactory.CreateTransactionScope(&suite.DataExecutorMock, func(tx data.Transaction) (bool, error) {
		suite.Equal(&suite.TransactionMock, tx)
		return false, nil
	})

	//assert
	suite.Require().NoError(err)

	suite.TransactionMock.AssertCalled(suite.T(), "Rollback")
	suite.TransactionMock.AssertNotCalled(suite.T(), "Commit")
}

func (suite *ScopeFactoryTestSuite) TestCreateTransactionScope_WithErrorCommitingTransaction_ReturnsError() {
	//arrange
	message := "commit error"

	suite.DataExecutorMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.TransactionMock.On("Rollback").Return(nil)
	suite.TransactionMock.On("Commit").Return(errors.New(message))

	//act
	err := suite.ScopeFactory.CreateTransactionScope(&suite.DataExecutorMock, func(tx data.Transaction) (bool, error) {
		suite.Equal(&suite.TransactionMock, tx)
		return true, nil
	})

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *ScopeFactoryTestSuite) TestCreateTransactionScope_WithSuccessFromBody_CommitsTransaction() {
	//arrange
	suite.DataExecutorMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.TransactionMock.On("Rollback").Return(nil)
	suite.TransactionMock.On("Commit").Return(nil)

	//act
	err := suite.ScopeFactory.CreateTransactionScope(&suite.DataExecutorMock, func(tx data.Transaction) (bool, error) {
		suite.Equal(&suite.TransactionMock, tx)
		return true, nil
	})

	//assert
	suite.Require().NoError(err)

	suite.DataExecutorMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.TransactionMock.AssertCalled(suite.T(), "Commit")
}

func TestScopeFactoryTestSuite(t *testing.T) {
	suite.Run(t, &ScopeFactoryTestSuite{})
}
