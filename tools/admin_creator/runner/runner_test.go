package runner_test

import (
	"authserver/common"
	controllermocks "authserver/controllers/mocks"
	"authserver/models"
	"authserver/testing/helpers"
	"authserver/tools/admin_creator/runner"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AdminCreatorTestSuite struct {
	suite.Suite
	helpers.ScopeFactorySuite
	ControllersMock controllermocks.Controllers
}

func (suite *AdminCreatorTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.ControllersMock = controllermocks.Controllers{}
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingUser_ReturnsError() {
	//arrange
	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.False(result)
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})

	//act
	err := runner.Run(&suite.ScopeFactoryMock, &suite.ControllersMock, "username", "password", 0)

	//assert
	suite.NoError(err)
}

func (suite *AdminCreatorTestSuite) TestRun_WithNoErrors_ReturnsNoErrors() {
	//arrange
	username := "username"
	password := "password"
	rank := 0

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, common.NoError())

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.True(result)
		suite.NoError(err)
	})

	//act
	err := runner.Run(&suite.ScopeFactoryMock, &suite.ControllersMock, username, password, rank)

	//assert
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateTransactionScope", &suite.DataExecutorMock, mock.Anything)
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.TransactionMock, username, password, rank)

	suite.NoError(err)
}

func TestAdminCreatorTestSuite(t *testing.T) {
	suite.Run(t, &AdminCreatorTestSuite{})
}
