package server_test

import (
	"authserver/server"
	"authserver/server/mocks"
	"authserver/testing/helpers"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	helpers.CustomSuite
	ServerMock mocks.Server
	Runner     *server.Runner
}

func (suite *RunnerTestSuite) SetupTest() {
	suite.ServerMock = mocks.Server{}

	suite.Runner = &server.Runner{
		Server: &suite.ServerMock,
	}
}

func (suite *RunnerTestSuite) TestRun_WithErrorStartingServer_ReturnsError() {
	//arrange
	message := "Start mock error"

	suite.ServerMock.On("Start").Return(errors.New(message))

	//act
	err := suite.Runner.Run()

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *RunnerTestSuite) TestRun_StartsServer() {
	//arrange
	suite.ServerMock.On("Start").Return(nil)

	//act
	err := suite.Runner.Run()

	//assert
	suite.Require().NoError(err)
	suite.ServerMock.AssertCalled(suite.T(), "Start")
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, &RunnerTestSuite{})
}
