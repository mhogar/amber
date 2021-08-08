package controllers_test

import (
	datamocks "authserver/data/mocks"

	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	CRUDMock datamocks.DataCRUD
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
}
