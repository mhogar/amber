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
	err := suite.Tx.SaveSession(models.CreateNewSession(nil))

	//assert
	suite.Require().Error(err)
	helpers.AssertContainsSubstrings(&suite.Suite, err.Error(), "error", "session model")
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_WhereSessionNotFound_ReturnsNilSession() {
	//act
	session, err := suite.Tx.GetSessionByID(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(session)
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_GetsTheSessionWithId() {
	//arrange
	session := models.CreateNewSession(
		models.CreateUser("username", []byte("password")),
	)
	suite.SaveSessionAndFields(session)

	//act
	resultSession, err := suite.Tx.GetSessionByID(session.ID)

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
	session := models.CreateNewSession(
		models.CreateUser("username", []byte("password")),
	)
	suite.SaveSessionAndFields(session)

	//act
	res, err := suite.Tx.DeleteSession(session.ID)
	suite.Require().NoError(err)

	//assert
	resultSession, err := suite.Tx.GetSessionByID(session.ID)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultSession)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_WithNoSessionsToDelete_ReturnsNilError() {
	//arrange
	session := models.CreateNewSession(
		models.CreateUser("", nil),
	)

	//act
	err := suite.Tx.DeleteAllOtherUserSessions(session.User.Username, session.ID)

	//assert
	suite.NoError(err)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_DeletesAllOtherSessionWithUserId() {
	//arrange
	session1 := models.CreateNewSession(
		models.CreateUser("username", []byte("password")),
	)
	suite.SaveSessionAndFields(session1)

	session2 := models.CreateNewSession(session1.User)
	suite.Tx.SaveSession(session2)

	//act
	err := suite.Tx.DeleteAllOtherUserSessions(session1.User.Username, session1.ID)

	//assert
	suite.Require().NoError(err)

	//can still find session1
	resultSession, err := suite.Tx.GetSessionByID(session1.ID)
	suite.NoError(err)
	suite.EqualValues(session1, resultSession)

	//session2 was deleted
	resultSession, err = suite.Tx.GetSessionByID(session2.ID)
	suite.NoError(err)
	suite.Nil(resultSession)
}

func TestSessionCRUDTestSuite(t *testing.T) {
	suite.Run(t, &SessionCRUDTestSuite{})
}
