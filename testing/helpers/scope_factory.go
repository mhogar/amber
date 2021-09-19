package helpers

import (
	"authserver/data"
	"authserver/data/mocks"

	"github.com/stretchr/testify/mock"
)

type ScopeFactorySuite struct {
	CustomSuite
	ScopeFactoryMock mocks.ScopeFactory
	DataExecutorMock mocks.DataExecutor
	TransactionMock  mocks.Transaction
}

func (suite *ScopeFactorySuite) SetupTest() {
	suite.ScopeFactoryMock = mocks.ScopeFactory{}
	suite.DataExecutorMock = mocks.DataExecutor{}
	suite.TransactionMock = mocks.Transaction{}

	suite.ScopeFactoryMock.On("CreateTransactionScope").Return(nil).Run(func(args mock.Arguments) {
		body := args.Get(1).(func(data.Transaction) (bool, error))
		body(&suite.TransactionMock)
	})
}

// SetupScopeFactoryMock_CreateDataExecutorScope sets up CreateDataExecutorScope with the provided result.
func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateDataExecutorScope(result error) {
	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(result, func(_ error) {})
}

// SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback sets up CreateDataExecutorScope with the provided result and callback.
func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(result error, callback func(error)) {
	suite.ScopeFactoryMock.On("CreateDataExecutorScope", mock.Anything).Return(result).Run(func(args mock.Arguments) {
		if result != nil {
			return
		}

		body := args.Get(0).(func(data.DataExecutor) error)
		err := body(&suite.DataExecutorMock)

		callback(err)
	})
}

// SetupScopeFactoryMock_CreateTransactionScope sets up CreateTransactionScope with the provided result.
func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateTransactionScope(result error) {
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(result, func(_ bool, _ error) {})
}

// SetupScopeFactoryMock_CreateTransactionScope_WithCallback sets up CreateTransactionScope with the provided result and callback.
func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateTransactionScope_WithCallback(result error, callback func(bool, error)) {
	suite.ScopeFactoryMock.On("CreateTransactionScope", mock.Anything, mock.Anything).Return(result).Run(func(args mock.Arguments) {
		if result != nil {
			return
		}

		body := args.Get(1).(func(data.Transaction) (bool, error))
		result, err := body(&suite.TransactionMock)

		callback(result, err)
	})
}
