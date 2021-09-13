package models

import "github.com/google/uuid"

const (
	ValidateUserRoleValid       = 0x0
	ValidateUserRoleEmptyRole   = 0x1
	ValidateUserRoleRoleTooLong = 0x2
)

// UserRoleRoleMaxLength is the max length a user's username can be.
const UserRoleRoleMaxLength = 15

// UserRole represents the user-role model
type UserRole struct {
	ClientUID uuid.UUID
	Username  string
	Role      string
}

type UserRoleCRUD interface {
	// CreateUserRole creates the user-role. Returns any errors.
	CreateUserRole(role *UserRole) error

	// GetUserRolesByClientUID fetches the user roles for the provided client uid and with a rank less than the provided rank.
	// Returns the user-roles and returns any errors.
	GetUserRolesWithLesserRankByClientUID(uid uuid.UUID, rank int) ([]*UserRole, error)

	// GetUserRoleByClientUIDAndUsername fetches the user role for the provided client uid and username.
	// Returns the user-role if it exists, nil if not. Also returns any errors.
	GetUserRoleByClientUIDAndUsername(clientUID uuid.UUID, username string) (*UserRole, error)

	// UpdateUserRole updates the user-role.
	// Returns result of whether the user was found and any errors.
	UpdateUserRole(role *UserRole) (bool, error)

	// DeleteUserRole deletes the user-role with the given username and client uid.
	// Returns result of whether the user-role was found, and any errors.
	DeleteUserRole(username string, clientUID uuid.UUID) (bool, error)
}

func CreateUserRole(clientUID uuid.UUID, username string, role string) *UserRole {
	return &UserRole{
		ClientUID: clientUID,
		Username:  username,
		Role:      role,
	}
}

// Validate validates the user-role model has valid fields.
// Returns an int indicating which fields are invalid.
func (ur *UserRole) Validate() int {
	code := ValidateUserRoleValid

	if ur.Role == "" {
		code |= ValidateUserRoleEmptyRole
	} else if len(ur.Role) > UserRoleRoleMaxLength {
		code |= ValidateUserRoleRoleTooLong
	}

	return code
}
