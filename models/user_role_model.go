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
	Username  string
	ClientUID uuid.UUID
	Role      string
}

type UserRoleCRUD interface {
	// CreateUserRole creates the user-role. Returns any errors.
	CreateUserRole(role *UserRole) error

	// GetUserRolesByClientUID fetches the user roles for the provided client uid.
	// Returns the user-roles and returns any errors.
	GetUserRolesByClientUID(uid uuid.UUID) ([]*UserRole, error)

	// GetUserRoleByUsernameAndClientUID fetches the user role for the provided username and client uid.
	// Returns the user-role if it exists, nil if not. Also returns any errors.
	GetUserRoleByUsernameAndClientUID(username string, clientUID uuid.UUID) (*UserRole, error)

	// UpdateUserRole updates the user-role.
	// Returns result of whether the user was found and any errors.
	UpdateUserRole(role *UserRole) (bool, error)

	// DeleteUserRole deletes the user-role with the given username and client uid.
	// Returns result of whether the user-role was found, and any errors.
	DeleteUserRole(username string, clientUID uuid.UUID) (bool, error)
}

func CreateUserRole(username string, clientUID uuid.UUID, role string) *UserRole {
	return &UserRole{
		Username:  username,
		ClientUID: clientUID,
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
