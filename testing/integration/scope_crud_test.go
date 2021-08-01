package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ScopeCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *ScopeCRUDTestSuite) TestSaveScope_WithInvalidScope_ReturnsError() {
	//arrange
	scope := models.CreateNewScope("")

	//act
	err := suite.Tx.SaveScope(scope)

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "scope model")
}

func (suite *ScopeCRUDTestSuite) TestGetScopeByName_WhereScopeNotFound_ReturnsNilScope() {
	//act
	scope, err := suite.Tx.GetScopeByName("not a real name")

	//assert
	suite.NoError(err)
	suite.Nil(scope)
}

func (suite *ScopeCRUDTestSuite) TestGetScopeByName_GetsTheScopeWithName() {
	//arrange
	scope := models.CreateNewScope("name")
	suite.SaveScope(suite.Tx, scope)

	//act
	resultScope, err := suite.Tx.GetScopeByName(scope.Name)

	//assert
	suite.NoError(err)
	suite.EqualValues(scope, resultScope)
}

func TestScopeCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ScopeCRUDTestSuite{})
}
