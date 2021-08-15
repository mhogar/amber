package models

const (
	ValidateUserValid               = 0x0
	ValidateUserEmptyUsername       = 0x1
	ValidateUserUsernameTooLong     = 0x2
	ValidateUserInvalidPasswordHash = 0x4
)

// UserUsernameMaxLength is the max length a user's username can be.
const UserUsernameMaxLength = 30

// User represents the user model.
type User struct {
	Username     string
	PasswordHash []byte
}

type UserCRUD interface {
	// CreateUser creates a new user and returns any errors.
	CreateUser(user *User) error

	// GetUserByUsername fetches the user with the matching username.
	// If no users are found, returns nil user. Also returns any errors.
	GetUserByUsername(username string) (*User, error)

	// UpdateUser updates the user and returns any errors.
	// Returns result of whether the user was found, and any errors.
	UpdateUser(user *User) (bool, error)

	// DeleteUser deletes the user with the given username.
	// Returns result of whether the user was found, and any errors.
	DeleteUser(username string) (bool, error)
}

func CreateUser(username string, passwordHash []byte) *User {
	return &User{
		Username:     username,
		PasswordHash: passwordHash,
	}
}

func CreateNewUser(username string, passwordHash []byte) *User {
	return CreateUser(username, passwordHash)
}

// Validate validates the user model has valid fields.
// Returns an int indicating which fields are invalid.
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
