package models

import "github.com/google/uuid"

const (
	ValidateSessionValid       = 0x0
	ValidateSessionNilID       = 0x1
	ValidateSessionNilUser     = 0x2
	ValidateSessionInvalidUser = 0x4
)

// Session represents the session model.
type Session struct {
	ID   uuid.UUID
	User *User
}

type SessionCRUD interface {
	// SaveSession saves the session and returns any errors.
	SaveSession(session *Session) error

	// GetSessionByID fetches the session associated with the id.
	// If no sessions are found, returns nil session.
	// Also returns any errors.
	GetSessionByID(ID uuid.UUID) (*Session, error)

	// DeleteSession deletes the session with the given id.
	// Returns result of whether the session was found, and any errors.
	DeleteSession(ID uuid.UUID) (bool, error)

	// DeleteAllOtherUserSessions deletes all of the sessions for the given username expect the one with the given id.
	// Returns any errors.
	DeleteAllOtherUserSessions(username string, ID uuid.UUID) error
}

func CreateSession(id uuid.UUID, user *User) *Session {
	return &Session{
		ID:   id,
		User: user,
	}
}

func CreateNewSession(user *User) *Session {
	return CreateSession(uuid.New(), user)
}

// Validate validates the access token model has valid fields.
// Returns an int indicating which fields are invalid.
func (tk *Session) Validate() int {
	code := ValidateSessionValid

	if tk.ID == uuid.Nil {
		code |= ValidateSessionNilID
	}

	if tk.User == nil {
		code |= ValidateSessionNilUser
	} else {
		verr := tk.User.Validate()
		if verr != ValidateUserValid {
			code |= ValidateSessionInvalidUser
		}
	}

	return code
}
