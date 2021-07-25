package main_test

import (
	requesterror "authserver/common/request_error"
	controllermocks "authserver/controllers/mocks"
	datamocks "authserver/data/mocks"
	"authserver/models"
	admincreator "authserver/tools/admin_creator"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AdminCreatorTestSuite struct {
	suite.Suite
	ControllersMock controllermocks.Controllers
	DataAdpaterMock datamocks.DataAdapter
	TransactionMock datamocks.Transaction
}

func (suite *AdminCreatorTestSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.DataAdpaterMock = datamocks.DataAdapter{}
	suite.TransactionMock = datamocks.Transaction{}

	suite.TransactionMock.On("Rollback")
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorSettingUpDataAdapter_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	message := "Setup test error"
	suite.DataAdpaterMock.On("Setup").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DataAdpaterMock, &suite.ControllersMock, username, password)

	//assert
	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingTransaction_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DataAdpaterMock.On("Setup").Return(nil)
	suite.DataAdpaterMock.On("CleanUp").Return(nil)

	message := "create transaction error"
	suite.DataAdpaterMock.On("CreateTransaction").Return(nil, errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DataAdpaterMock, &suite.ControllersMock, username, password)

	//assert
	suite.DataAdpaterMock.AssertCalled(suite.T(), "CleanUp")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingUser_RollbacksTransactionReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DataAdpaterMock.On("Setup").Return(nil)
	suite.DataAdpaterMock.On("CleanUp").Return(nil)
	suite.DataAdpaterMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	user, err := admincreator.Run(&suite.DataAdpaterMock, &suite.ControllersMock, username, password)

	//assert
	suite.DataAdpaterMock.AssertCalled(suite.T(), "CleanUp")
	suite.TransactionMock.AssertCalled(suite.T(), "Rollback")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCommitingTransaction_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DataAdpaterMock.On("Setup").Return(nil)
	suite.DataAdpaterMock.On("CleanUp").Return(nil)
	suite.ControllersMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())

	message := "commit transaction error"
	suite.TransactionMock.On("Commit").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DataAdpaterMock, &suite.ControllersMock, username, password)

	//assert
	suite.DataAdpaterMock.AssertCalled(suite.T(), "CleanUp")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithNoErrors_ReturnsNoErrors() {
	//arrange
	username := "username"
	password := "password"

	suite.DataAdpaterMock.On("Setup").Return(nil)
	suite.DataAdpaterMock.On("CleanUp").Return(nil)
	suite.DataAdpaterMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	user, err := admincreator.Run(&suite.DataAdpaterMock, &suite.ControllersMock, username, password)

	//assert
	suite.DataAdpaterMock.AssertCalled(suite.T(), "Setup")
	suite.DataAdpaterMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", mock.Anything, username, password)
	suite.TransactionMock.AssertCalled(suite.T(), "Commit")
	suite.TransactionMock.AssertNotCalled(suite.T(), "Rollback")
	suite.DataAdpaterMock.AssertCalled(suite.T(), "Setup")

	suite.NoError(err)
	suite.NotNil(user)
}

func TestAdminCreatorTestSuite(t *testing.T) {
	suite.Run(t, &AdminCreatorTestSuite{})
}
