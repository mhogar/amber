package firestoreadapter

import (
	"errors"
	"fmt"

	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
)

func (crud *FirestoreCRUD) CreateUserRole(role *models.UserRole) error {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	return nil
}

func (crud *FirestoreCRUD) GetUserRolesWithLesserRankByClientUID(uid uuid.UUID, rank int) ([]*models.UserRole, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) GetUserRoleByClientUIDAndUsername(clientUID uuid.UUID, username string) (*models.UserRole, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) UpdateUserRole(role *models.UserRole) (bool, error) {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return false, errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	return false, nil
}

func (crud *FirestoreCRUD) DeleteUserRole(username string, clientUID uuid.UUID) (bool, error) {
	return false, nil
}
