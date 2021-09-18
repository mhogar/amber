package server_test

import (
	routermocks "authserver/router/mocks"
	"authserver/server"
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	helpers.CustomSuite
	RouterFactoryMock routermocks.RouterFactory
}

func (suite *ServerTestSuite) SetupTest() {
	suite.RouterFactoryMock = routermocks.RouterFactory{}
}

func (suite *ServerTestSuite) TestCreateHTTPServerRunner_CreatesRunnerUsingHTTPServer() {
	//arrange
	suite.RouterFactoryMock.On("CreateRouter").Return(nil)

	//act
	runner := server.CreateHTTPServerRunner(&suite.RouterFactoryMock)
	_, ok := runner.Server.(*server.HTTPServer)

	//assert
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
	suite.True(ok, "Runner's server should be an http server")
}

func (suite *ServerTestSuite) TestCreateHTTPTestServerRunner_CreatesRunnerUsingHTTPTestServer() {
	//arrange
	suite.RouterFactoryMock.On("CreateRouter").Return(nil)

	//act
	runner := server.CreateHTTPTestServerRunner(&suite.RouterFactoryMock)
	_, ok := runner.Server.(*server.HTTPTestServer)

	//assert
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
	suite.True(ok, "Runner's server should be an httptest server")
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}
