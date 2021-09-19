package controllers_test

import (
	datamocks "github.com/mhogar/amber/data/mocks"
	"github.com/mhogar/amber/testing/helpers"
)

type ControllerTestSuite struct {
	helpers.CustomSuite
	CRUDMock datamocks.DataCRUD
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
}
