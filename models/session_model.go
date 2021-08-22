package models

import "github.com/google/uuid"

const (
	ValidateSessionValid         = 0x0
	ValidateSessionNilToken      = 0x1
	ValidateSessionEmptyUsername = 0x2
)

// Session represents the session model.
type Session struct {
	Token    uuid.UUID
	Username string
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

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username expect the one with the given token.
	// Returns any errors.
	DeleteAllOtherUserSessions(username string, tokem uuid.UUID) error
}

func CreateSession(token uuid.UUID, username string) *Session {
	return &Session{
		Token:    token,
		Username: username,
	}
}

func CreateNewSession(username string) *Session {
	return CreateSession(uuid.New(), username)
}

// Validate validates the access token model has valid fields.
// Returns an int indicating which fields are invalid.
func (tk *Session) Validate() int {
	code := ValidateSessionValid

	if tk.Token == uuid.Nil {
		code |= ValidateSessionNilToken
	}

	if tk.Username == "" {
		code |= ValidateSessionEmptyUsername
	}

	return code
}
