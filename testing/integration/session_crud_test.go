package integration_test

import (
	"authserver/models"
	"authserver/testing/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SessionCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *SessionCRUDTestSuite) TestSaveSession_WithInvalidSession_ReturnsError() {
	//act
	err := suite.Tx.SaveSession(models.CreateSession(uuid.Nil, ""))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "session model")
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_WhereSessionNotFound_ReturnsNilSession() {
	//act
	session, err := suite.Tx.GetSessionByToken(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(session)
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_GetsTheSessionWithId() {
	//arrange
	user := models.CreateUser("username", []byte("password"))
	suite.SaveUser(user)

	session := models.CreateNewSession(user.Username)
	suite.SaveSession(session)

	//act
	resultSession, err := suite.Tx.GetSessionByToken(session.Token)

	//assert
	suite.NoError(err)
	suite.EqualValues(session, resultSession)
}

func (suite *SessionCRUDTestSuite) TestDeleteSession_WhereSessionIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Tx.DeleteSession(uuid.New())

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *SessionCRUDTestSuite) TestDeleteSession_DeletesSessionWithId() {
	//arrange
	user := models.CreateUser("username", []byte("password"))
	suite.SaveUser(user)

	session := models.CreateNewSession(user.Username)
	suite.SaveSession(session)

	//act
	res, err := suite.Tx.DeleteSession(session.Token)
	suite.Require().NoError(err)

	//assert
	resultSession, err := suite.Tx.GetSessionByToken(session.Token)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultSession)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_WithNoSessionsToDelete_ReturnsNilError() {
	//arrange
	user := models.CreateUser("username", []byte("password"))
	suite.SaveUser(user)

	session := models.CreateNewSession(user.Username)
	suite.SaveSession(session)

	//act
	err := suite.Tx.DeleteAllOtherUserSessions(session.Username, session.Token)

	//assert
	suite.NoError(err)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_DeletesAllOtherSessionWithUserId() {
	//arrange
	user := models.CreateUser("username", []byte("password"))
	suite.SaveUser(user)

	session1 := models.CreateNewSession(user.Username)
	suite.SaveSession(session1)

	session2 := models.CreateNewSession(session1.Username)
	suite.Tx.SaveSession(session2)

	//act
	err := suite.Tx.DeleteAllOtherUserSessions(session1.Username, session1.Token)

	//assert
	suite.Require().NoError(err)

	//can still find session1
	resultSession, err := suite.Tx.GetSessionByToken(session1.Token)
	suite.NoError(err)
	suite.EqualValues(session1, resultSession)

	//session2 was deleted
	resultSession, err = suite.Tx.GetSessionByToken(session2.Token)
	suite.NoError(err)
	suite.Nil(resultSession)
}

func TestSessionCRUDTestSuite(t *testing.T) {
	suite.Run(t, &SessionCRUDTestSuite{})
}
