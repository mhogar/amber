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

	//assert
	suite.IsType(&server.HTTPServer{}, runner.Server)
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
}

func (suite *ServerTestSuite) TestCreateHTTPTestServerRunner_CreatesRunnerUsingHTTPTestServer() {
	//arrange
	suite.RouterFactoryMock.On("CreateRouter").Return(nil)

	//act
	runner := server.CreateHTTPTestServerRunner(&suite.RouterFactoryMock)

	//assert
	suite.IsType(&server.HTTPTestServer{}, runner.Server)
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}
