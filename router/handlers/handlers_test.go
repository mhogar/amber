package handlers_test

import (
	controllermocks "authserver/controllers/mocks"
	datamocks "authserver/data/mocks"
	"authserver/router/handlers"

	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	ControllersMock controllermocks.Controllers
	TransactionMock datamocks.Transaction
	CoreHandlers    handlers.CoreHandlers
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.TransactionMock = datamocks.Transaction{}

	suite.CoreHandlers = handlers.CoreHandlers{
		Controllers: &suite.ControllersMock,
	}
}
