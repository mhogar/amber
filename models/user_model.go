package models

import (
	"github.com/google/uuid"
)

const (
	ValidateUserValid               = 0x0
	ValidateUserEmptyUsername       = 0x1
	ValidateUserUsernameTooLong     = 0x2
	ValidateUserInvalidPasswordHash = 0x4
)

// UserUsernameMaxLength is the max length a user's username can be
const UserUsernameMaxLength = 30

// User represents the user model
type User struct {
	ID           int32
	Username     string
	PasswordHash []byte
}

type UserCRUD interface {
	// SaveUser saves the user and returns any errors
	SaveUser(user *User) error

	// GetUserByID fetches the user associated with the id
	// If no users are found, returns nil user. Also returns any errors
	GetUserByID(ID uuid.UUID) (*User, error)

	// GetUserByUsername fetches the user with the matching username
	// If no users are found, returns nil user. Also returns any errors
	GetUserByUsername(username string) (*User, error)

	// UpdateUser updates the user and returns any errors
	UpdateUser(user *User) error

	// DeleteUser deletes the user and returns any errors
	DeleteUser(user *User) error
}

func CreateUser(id int32, username string, passwordHash []byte) *User {
	return &User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
	}
}

func CreateNewUser(username string, passwordHash []byte) *User {
	return CreateUser(0, username, passwordHash)
}

// Validate validates the user model has valid fields
// Returns an int indicating which fields are invalid
func (u *User) Validate() int {
	code := ValidateUserValid

	if u.Username == "" {
		code |= ValidateUserEmptyUsername
	} else if len(u.Username) > UserUsernameMaxLength {
		code |= ValidateUserUsernameTooLong
	}

	if len(u.PasswordHash) == 0 {
		code |= ValidateUserInvalidPasswordHash
	}

	return code
}
