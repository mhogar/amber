package testhelpers

import (
	"authserver/data"
	"authserver/data/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ScopeFactorySuite struct {
	suite.Suite
	ScopeFactoryMock mocks.IScopeFactory
	DataExecutorMock mocks.DataExecutor
	TransactionMock  mocks.Transaction
}

func (suite *ScopeFactorySuite) SetupTest() {
	suite.ScopeFactoryMock = mocks.IScopeFactory{}
	suite.DataExecutorMock = mocks.DataExecutor{}
	suite.TransactionMock = mocks.Transaction{}

	suite.ScopeFactoryMock.On("CreateTransactionScope").Return(nil).Run(func(args mock.Arguments) {
		body := args.Get(1).(func(data.Transaction) (bool, error))
		body(&suite.TransactionMock)
	})
}

func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateDataExecutorScope(result error) {
	suite.SetupScopeFactoryMock_CreateDataExecutorScope_WithCallback(result, func(_ error) {})
}

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

func (suite *ScopeFactorySuite) SetupScopeFactoryMock_CreateTransactionScope(result error) {
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(result, func(_ bool, _ error) {})
}

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
