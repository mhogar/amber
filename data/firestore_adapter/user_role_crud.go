package firestoreadapter

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"google.golang.org/api/iterator"

	"github.com/google/uuid"
)

func (crud *FirestoreCRUD) CreateUserRole(role *models.UserRole) error {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	//create user-role
	err := crud.DocWriter.Create(crud.getUserRoleDocRef(role.ClientUID, role.Username), role)
	if err != nil {
		return common.ChainError("error creating user-role", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetUserRolesWithLesserRankByClientUID(uid uuid.UUID, rank int) ([]*models.UserRole, error) {
	//get users
	users, err := crud.GetUsersWithLesserRank(rank)
	if err != nil {
		return nil, common.ChainError("error getting user with lesser rank", err)
	}

	//extract usernames
	usernames := make([]string, len(users))
	for index, user := range users {
		usernames[index] = user.Username
	}

	//get user-roles
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("user-roles").
		Where("client_uid", "==", uid).
		Where("username", "in", usernames).
		OrderBy("username", firestore.Asc).
		Documents(ctx)

	defer cancel()
	defer itr.Stop()

	//read the results
	roles := []*models.UserRole{}
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, common.ChainError("error getting next doc", err)
		}

		role, err := crud.readUserRoleData(doc)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (crud *FirestoreCRUD) GetUserRoleByClientUIDAndUsername(clientUID uuid.UUID, username string) (*models.UserRole, error) {
	doc, err := crud.getUserRole(clientUID, username)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	return crud.readUserRoleData(doc)
}

func (crud *FirestoreCRUD) UpdateUserRole(role *models.UserRole) (bool, error) {
	//validate the user-role model
	verr := role.Validate()
	if verr != models.ValidateUserRoleValid {
		return false, errors.New(fmt.Sprint("error validating user-role model:", verr))
	}

	//check user-role already exists
	doc, err := crud.getUserRole(role.ClientUID, role.Username)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//update user-role
	err = crud.DocWriter.Set(doc.Ref, role)
	if err != nil {
		return true, common.ChainError("error updating user-role", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) DeleteUserRole(clientUID uuid.UUID, username string) (bool, error) {
	//check user-role already exists
	doc, err := crud.getUserRole(clientUID, username)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//delete user-role
	err = crud.DocWriter.Delete(doc.Ref)
	if err != nil {
		return false, common.ChainError("error deleting user-role", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) DeleteAllUserRolesByClientUID(uid uuid.UUID) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("user-roles").
		Where("client_uid", "==", uid).
		Documents(ctx)
	defer cancel()

	//delete user-roles
	err := crud.deleteUserRoles(itr)
	if err != nil {
		return err
	}

	return nil
}

func (crud *FirestoreCRUD) DeleteAllUserRolesByUsername(username string) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("user-roles").
		Where("username", "==", username).
		Documents(ctx)
	defer cancel()

	//delete user-roles
	err := crud.deleteUserRoles(itr)
	if err != nil {
		return err
	}

	return nil
}

func (crud *FirestoreCRUD) deleteUserRoles(itr *firestore.DocumentIterator) error {
	defer itr.Stop()
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return common.ChainError("error getting next doc", err)
		}

		//delete user-role
		err = crud.DocWriter.Delete(doc.Ref)
		if err != nil {
			return common.ChainError("error deleting user-role", err)
		}
	}
}

func (crud *FirestoreCRUD) getUserRoleDocRef(clientUID uuid.UUID, username string) *firestore.DocumentRef {
	return crud.Client.Collection("user-roles").Doc(clientUID.String() + "-" + username)
}

func (crud *FirestoreCRUD) getUserRole(clientUID uuid.UUID, username string) (*firestore.DocumentSnapshot, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	doc, err := crud.getUserRoleDocRef(clientUID, username).Get(ctx)
	cancel()

	//check user-role was found
	if !doc.Exists() {
		return nil, nil
	}

	//handle other errors
	if err != nil {
		return nil, common.ChainError("error getting user-role", err)
	}

	return doc, nil
}

func (*FirestoreCRUD) readUserRoleData(doc *firestore.DocumentSnapshot) (*models.UserRole, error) {
	role := &models.UserRole{}

	err := doc.DataTo(&role)
	if err != nil {
		return nil, common.ChainError("error reading user-role data", err)
	}

	return role, nil
}
