package main_test

import (
	requesterror "authserver/common/request_error"
	controllermocks "authserver/controllers/mocks"
	"authserver/models"
	testhelpers "authserver/testing"
	admincreator "authserver/tools/admin_creator"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AdminCreatorTestSuite struct {
	suite.Suite
	testhelpers.ScopeFactorySuite
	ControllersMock controllermocks.Controllers
}

func (suite *AdminCreatorTestSuite) SetupTest() {
	suite.ScopeFactorySuite.SetupTest()
	suite.ControllersMock = controllermocks.Controllers{}
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingUser_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.False(result)
		suite.Require().Error(err)
		suite.Contains(err.Error(), message)
	})

	//act
	err := admincreator.Run(&suite.ScopeFactoryMock, &suite.ControllersMock, username, password)

	//assert
	suite.NoError(err)
}

func (suite *AdminCreatorTestSuite) TestRun_WithNoErrors_ReturnsNoErrors() {
	//arrange
	username := "username"
	password := "password"

	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())

	suite.SetupScopeFactoryMock_CreateDataExecutorScope(nil)
	suite.SetupScopeFactoryMock_CreateTransactionScope_WithCallback(nil, func(result bool, err error) {
		suite.True(result)
		suite.NoError(err)
	})

	//act
	err := admincreator.Run(&suite.ScopeFactoryMock, &suite.ControllersMock, username, password)

	//assert
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateDataExecutorScope", mock.Anything)
	suite.ScopeFactoryMock.AssertCalled(suite.T(), "CreateTransactionScope", &suite.DataExecutorMock, mock.Anything)
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", &suite.TransactionMock, username, password)

	suite.NoError(err)
}

func TestAdminCreatorTestSuite(t *testing.T) {
	suite.Run(t, &AdminCreatorTestSuite{})
}
