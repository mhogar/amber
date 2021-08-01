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
	Handlers        handlers.Handlers
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.TransactionMock = datamocks.Transaction{}

	suite.Handlers = handlers.Handlers{
		Controllers: &suite.ControllersMock,
	}
}
