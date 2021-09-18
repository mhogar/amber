package controllers_test

import (
	datamocks "authserver/data/mocks"
	"authserver/testing/helpers"
)

type ControllerTestSuite struct {
	helpers.CustomSuite
	CRUDMock datamocks.DataCRUD
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
}
