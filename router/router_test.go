package router_test

import (
	databasemocks "authserver/database/mocks"
	"authserver/router"
	"authserver/router/mocks"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
	AuthenticatorMock      mocks.Authenticator
	TransactionFactoryMock databasemocks.TransactionFactory
	TransactionMock        databasemocks.Transaction
	Router                 *httprouter.Router
}

func (suite *RouterTestSuite) SetupTest() {
	suite.AuthenticatorMock = mocks.Authenticator{}
	suite.TransactionFactoryMock = databasemocks.TransactionFactory{}
	suite.TransactionMock = databasemocks.Transaction{}

	suite.TransactionMock.On("RollbackTransaction")

	rf := router.RouterFactory{
		Authenticator:      &suite.AuthenticatorMock,
		TransactionFactory: &suite.TransactionFactoryMock,
	}
	suite.Router = rf.CreateRouter()
}
