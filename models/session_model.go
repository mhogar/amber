package models

import "github.com/google/uuid"

const (
	ValidateSessionValid    = 0x0
	ValidateSessionNilToken = 0x1
)

// Session represents the session model.
type Session struct {
	Token    uuid.UUID `firestore:"token"`
	Username string    `firestore:"username"`
	Rank     int       `firestore:"rank"`
}

type SessionCRUD interface {
	// SaveSession saves the session and returns any errors.
	SaveSession(session *Session) error

	// GetSessionByToken fetches the session with the given token.
	// If no sessions are found, returns nil session.
	// Also returns any errors.
	GetSessionByToken(token uuid.UUID) (*Session, error)

	// DeleteSession deletes the session with the given token.
	// Returns result of whether the session was found, and any errors.
	DeleteSession(token uuid.UUID) (bool, error)

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username.
	// Returns any errors.
	DeleteAllUserSessions(username string) error

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username except the one with the given token.
	// Returns any errors.
	DeleteAllOtherUserSessions(username string, tokem uuid.UUID) error
}

// CreateSession creates a new session model with the provided field.
func CreateSession(token uuid.UUID, username string, rank int) *Session {
	return &Session{
		Token:    token,
		Username: username,
		Rank:     rank,
	}
}

// CreateNewSession generates a new token then creates a new session model with the token and provided field.
func CreateNewSession(username string, rank int) *Session {
	return CreateSession(uuid.New(), username, rank)
}

// Validate validates the access token model has valid fields.
// Returns an int indicating which fields are invalid.
func (s *Session) Validate() int {
	code := ValidateSessionValid

	//validate token
	if s.Token == uuid.Nil {
		code |= ValidateSessionNilToken
	}

	return code
}
