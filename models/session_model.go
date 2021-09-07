package models

import "github.com/google/uuid"

const (
	ValidateSessionValid           = 0x0
	ValidateSessionNilToken        = 0x1
	ValidateSessionEmptyUsername   = 0x2
	ValidateSessionUsernameTooLong = 0x4
	ValidateSessionInvalidRank     = 0x8
)

// Session represents the session model.
type Session struct {
	Token    uuid.UUID
	Username string
	Rank     int
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

func CreateSession(token uuid.UUID, username string, rank int) *Session {
	return &Session{
		Token:    token,
		Username: username,
		Rank:     rank,
	}
}

func CreateNewSession(username string, rank int) *Session {
	return CreateSession(uuid.New(), username, rank)
}

// Validate validates the access token model has valid fields.
// Returns an int indicating which fields are invalid.
func (s *Session) Validate() int {
	code := ValidateSessionValid

	if s.Token == uuid.Nil {
		code |= ValidateSessionNilToken
	}

	if s.Username == "" {
		code |= ValidateSessionEmptyUsername
	} else if len(s.Username) > UserUsernameMaxLength {
		code |= ValidateSessionUsernameTooLong
	}

	if s.Rank < 0 {
		code |= ValidateSessionInvalidRank
	}

	return code
}
