package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccessTokenCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *AccessTokenCRUDTestSuite) TestSaveAccessToken_WithInvalidAccessToken_ReturnsError() {
	//act
	err := suite.Tx.SaveAccessToken(models.CreateNewAccessToken(nil, nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "access token model")
}

func (suite *AccessTokenCRUDTestSuite) TestGetAccessTokenById_WhereAccessTokenNotFound_ReturnsNilAccessToken() {
	//act
	token, err := suite.Tx.GetAccessTokenByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(token)
}

func (suite *AccessTokenCRUDTestSuite) TestGetAccessTokenById_GetsTheAccessTokenWithId() {
	//arrange
	token := models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient("name"),
	)
	suite.SaveAccessTokenAndFields(token)

	//act
	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)

	//assert
	suite.NoError(err)
	suite.EqualValues(token, resultAccessToken)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAccessToken_WithNoAccessTokenToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteAccessToken(models.CreateNewAccessToken(nil, nil))

	//assert
	suite.NoError(err)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAccessToken_DeletesAccessTokenWithId() {
	//arrange
	token := models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient("name"),
	)
	suite.SaveAccessTokenAndFields(token)

	//act
	err := suite.Tx.DeleteAccessToken(token)

	//assert
	suite.Require().NoError(err)

	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token.ID)
	suite.NoError(err)
	suite.Nil(resultAccessToken)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAllOtherUserTokens_WithNoAccessTokensToDelete_ReturnsNilError() {
	//act
	err := suite.Tx.DeleteAllOtherUserTokens(models.CreateNewAccessToken(models.CreateNewUser("", nil), nil))

	//assert
	suite.NoError(err)
}

func (suite *AccessTokenCRUDTestSuite) TestDeleteAllOtherUserTokens_DeletesAllOtherTokenWithUserId() {
	//arrange
	token1 := models.CreateNewAccessToken(
		models.CreateNewUser("username", []byte("password")),
		models.CreateNewClient("name"),
	)
	suite.SaveAccessTokenAndFields(token1)

	token2 := models.CreateNewAccessToken(
		token1.User,
		token1.Client,
	)
	suite.Tx.SaveAccessToken(token2)

	//act
	err := suite.Tx.DeleteAllOtherUserTokens(token1)

	//assert
	suite.Require().NoError(err)

	//can still find token1
	resultAccessToken, err := suite.Tx.GetAccessTokenByID(token1.ID)
	suite.NoError(err)
	suite.EqualValues(token1, resultAccessToken)

	//token2 was deleted
	resultAccessToken, err = suite.Tx.GetAccessTokenByID(token2.ID)
	suite.NoError(err)
	suite.Nil(resultAccessToken)
}

func TestAccessTokenCRUDTestSuite(t *testing.T) {
	suite.Run(t, &AccessTokenCRUDTestSuite{})
}
