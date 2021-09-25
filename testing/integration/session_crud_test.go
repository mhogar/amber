package integration_test

import (
	"testing"

	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SessionCRUDTestSuite struct {
	CRUDTestSuite
}

func (suite *SessionCRUDTestSuite) TestSaveSession_WithInvalidSession_ReturnsError() {
	//act
	err := suite.Executor.SaveSession(models.CreateSession(uuid.Nil, "", -1))

	//assert
	suite.Require().Error(err)
	suite.ContainsSubstrings(err.Error(), "error", "session model")
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_WhereSessionNotFound_ReturnsNilSession() {
	//act
	session, err := suite.Executor.GetSessionByToken(uuid.New())

	//assert
	suite.NoError(err)
	suite.Nil(session)
}

func (suite *SessionCRUDTestSuite) TestGetSessionById_GetsTheSessionWithId() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	session := suite.SaveSession(models.CreateNewSession(user.Username, 0))

	//act
	resultSession, err := suite.Executor.GetSessionByToken(session.Token)

	//assert
	suite.NoError(err)
	suite.EqualValues(session, resultSession)

	//clean up
	suite.DeleteUser(user)
}

func (suite *SessionCRUDTestSuite) TestDeleteSession_WhereSessionIsNotFound_ReturnsFalseResult() {
	//act
	res, err := suite.Executor.DeleteSession(uuid.New())

	//assert
	suite.False(res)
	suite.NoError(err)
}

func (suite *SessionCRUDTestSuite) TestDeleteSession_DeletesSessionWithId() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	session := suite.SaveSession(models.CreateNewSession(user.Username, 0))

	//act
	res, err := suite.Executor.DeleteSession(session.Token)
	suite.Require().NoError(err)

	//assert
	resultSession, err := suite.Executor.GetSessionByToken(session.Token)

	suite.True(res)
	suite.NoError(err)
	suite.Nil(resultSession)

	//clean up
	suite.DeleteUser(user)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllUserSessions_WithNoSessionsToDelete_ReturnsNilError() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))

	//act
	err := suite.Executor.DeleteAllUserSessions(user.Username)

	//assert
	suite.NoError(err)

	//clean up
	suite.DeleteUser(user)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllUserSessions_DeletesAllSessionsWithUsername() {
	//arrange
	user1 := suite.SaveUser(models.CreateUser("user1", 0, []byte("password")))
	user2 := suite.SaveUser(models.CreateUser("user2", 0, []byte("password")))

	session1 := suite.SaveSession(models.CreateNewSession(user1.Username, 0))
	session2 := suite.SaveSession(models.CreateNewSession(user1.Username, 0))
	session3 := suite.SaveSession(models.CreateNewSession(user2.Username, 0))

	//act
	err := suite.Executor.DeleteAllUserSessions(user1.Username)

	//assert
	suite.Require().NoError(err)

	//session1 was deleted
	resultSession, err := suite.Executor.GetSessionByToken(session1.Token)
	suite.NoError(err)
	suite.Nil(resultSession)

	//session2 was deleted
	resultSession, err = suite.Executor.GetSessionByToken(session2.Token)
	suite.NoError(err)
	suite.Nil(resultSession)

	//can still find session3
	resultSession, err = suite.Executor.GetSessionByToken(session3.Token)
	suite.NoError(err)
	suite.EqualValues(session3, resultSession)

	//clean up
	suite.DeleteUser(user1)
	suite.DeleteUser(user2)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_WithNoSessionsToDelete_ReturnsNilError() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	session := suite.SaveSession(models.CreateNewSession(user.Username, 0))

	//act
	err := suite.Executor.DeleteAllOtherUserSessions(user.Username, session.Token)

	//assert
	suite.NoError(err)

	//clean up
	suite.DeleteUser(user)
}

func (suite *SessionCRUDTestSuite) TestDeleteAllOtherUserSessions_DeletesAllOtherSessionWithUsername() {
	//arrange
	user := suite.SaveUser(models.CreateUser("username", 0, []byte("password")))
	session1 := suite.SaveSession(models.CreateNewSession(user.Username, 0))
	session2 := suite.SaveSession(models.CreateNewSession(user.Username, 0))

	//act
	err := suite.Executor.DeleteAllOtherUserSessions(user.Username, session1.Token)

	//assert
	suite.Require().NoError(err)

	//can still find session1
	resultSession, err := suite.Executor.GetSessionByToken(session1.Token)
	suite.NoError(err)
	suite.EqualValues(session1, resultSession)

	//session2 was deleted
	resultSession, err = suite.Executor.GetSessionByToken(session2.Token)
	suite.NoError(err)
	suite.Nil(resultSession)

	//clean up
	suite.DeleteUser(user)
}

func TestSessionCRUDTestSuite(t *testing.T) {
	suite.Run(t, &SessionCRUDTestSuite{})
}
