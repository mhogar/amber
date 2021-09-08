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
	// GetUserRolesForClient fetches all the user roles for the provided client uid.
	// Returns a slice of the cleint's user-roles if they exist, nil if not.
	// Also returns any errors.
	GetUserRolesForClient(clientUID uuid.UUID) ([]*UserRole, error)

	// GetUserRoleForClientAndUser fetches the user roles for the provided client uid and username.
	// Returns the user-role if it exists, nil if not.
	// Also returns any errors.
	GetUserRoleForClient(clientUID uuid.UUID, username string) (*UserRole, error)

	// UpdateUserRolesForClient updates the roles for the provided client uid.
	// Returns any errors.
	UpdateUserRolesForClient(clientUID uuid.UUID, roles []*UserRole) error
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
