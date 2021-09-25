package firestoreadapter

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"google.golang.org/api/iterator"
)

func (crud *FirestoreCRUD) CreateUser(user *models.User) error {
	//validate the user model
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return errors.New(fmt.Sprint("error validating user model:", verr))
	}

	if user.PasswordHash == nil {
		return errors.New("password hash cannot be nil")
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	err := crud.DocWriter.Create(ctx, crud.Client.Collection("users").Doc(user.Username), user)
	cancel()

	if err != nil {
		return common.ChainError("error creating user", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetUsersWithLesserRank(rank int) ([]*models.User, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("users").
		Where("rank", "<", rank).
		OrderBy("rank", firestore.Asc).
		OrderBy("username", firestore.Asc).Documents(ctx)

	defer cancel()
	defer itr.Stop()

	//read the results
	users := []*models.User{}
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, common.ChainError("error getting next doc", err)
		}

		user, err := crud.readUserData(doc)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (crud *FirestoreCRUD) GetUserByUsername(username string) (*models.User, error) {
	doc, err := crud.getUser(username)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	return crud.readUserData(doc)
}

func (crud *FirestoreCRUD) UpdateUser(user *models.User) (bool, error) {
	//validate the user model
	verr := user.Validate()
	if verr != models.ValidateUserValid {
		return false, errors.New(fmt.Sprint("error validating user model:", verr))
	}

	//check user already exists
	doc, err := crud.getUser(user.Username)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//update fields
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	err = crud.DocWriter.Update(ctx, doc.Ref, []firestore.Update{
		{Path: "username", Value: user.Username},
		{Path: "rank", Value: user.Rank},
	})
	cancel()

	if err != nil {
		return true, common.ChainError("error updating user", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) UpdateUserPassword(username string, hash []byte) (bool, error) {
	if hash == nil {
		return false, errors.New("password hash cannot be nil")
	}

	//check user already exists
	doc, err := crud.getUser(username)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//update fields
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	err = crud.DocWriter.Update(ctx, doc.Ref, []firestore.Update{
		{Path: "password_hash", Value: hash},
	})
	cancel()

	if err != nil {
		return false, common.ChainError("error updating user password", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) DeleteUser(username string) (bool, error) {
	//check user already exists
	doc, err := crud.getUser(username)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//delete user
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	err = crud.DocWriter.Delete(ctx, doc.Ref)
	cancel()

	if err != nil {
		return false, common.ChainError("error deleting user", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) getUser(username string) (*firestore.DocumentSnapshot, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	doc, err := crud.Client.Collection("users").Doc(username).Get(ctx)
	cancel()

	//check user was found
	if !doc.Exists() {
		return nil, nil
	}

	//handle other errors
	if err != nil {
		return nil, common.ChainError("error getting user", err)
	}

	return doc, nil
}

func (*FirestoreCRUD) readUserData(doc *firestore.DocumentSnapshot) (*models.User, error) {
	user := &models.User{}

	err := doc.DataTo(&user)
	if err != nil {
		return nil, common.ChainError("error reading user data", err)
	}

	return user, nil
}
