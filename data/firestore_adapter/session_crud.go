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

func (crud *FirestoreCRUD) SaveSession(session *models.Session) error {
	//validate the session model
	verr := session.Validate()
	if verr != models.ValidateSessionValid {
		return errors.New(fmt.Sprint("error validating session model:", verr))
	}

	//create session
	err := crud.DocWriter.Create(crud.Client.Collection("sessions").Doc(session.Token.String()), session)
	if err != nil {
		return common.ChainError("error creating session", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	doc, err := crud.getSession(token)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	return crud.readSessionData(doc)
}

func (crud *FirestoreCRUD) DeleteSession(token uuid.UUID) (bool, error) {
	//check session already exists
	doc, err := crud.getSession(token)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//delete session
	err = crud.DocWriter.Delete(doc.Ref)
	if err != nil {
		return false, common.ChainError("error deleting session", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) DeleteAllUserSessions(username string) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("sessions").
		Where("username", "==", username).
		Documents(ctx)
	defer cancel()

	//delete the sessions
	crud.DeleteSessions(itr)

	return nil
}

func (crud *FirestoreCRUD) DeleteAllOtherUserSessions(username string, token uuid.UUID) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("sessions").
		Where("username", "==", username).
		Where("token", "!=", token).
		Documents(ctx)
	defer cancel()

	//delete the sessions
	crud.DeleteSessions(itr)

	return nil
}

func (crud *FirestoreCRUD) getSession(token uuid.UUID) (*firestore.DocumentSnapshot, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	doc, err := crud.Client.Collection("sessions").Doc(token.String()).Get(ctx)
	cancel()

	//check session was found
	if !doc.Exists() {
		return nil, nil
	}

	//handle other errors
	if err != nil {
		return nil, common.ChainError("error getting session", err)
	}

	return doc, nil
}

func (crud *FirestoreCRUD) DeleteSessions(itr *firestore.DocumentIterator) error {
	defer itr.Stop()
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return common.ChainError("error getting next doc", err)
		}

		//delete session
		err = crud.DocWriter.Delete(doc.Ref)
		if err != nil {
			return common.ChainError("error deleting session", err)
		}
	}
}

func (*FirestoreCRUD) readSessionData(doc *firestore.DocumentSnapshot) (*models.Session, error) {
	session := &models.Session{}

	err := doc.DataTo(&session)
	if err != nil {
		return nil, common.ChainError("error reading session data", err)
	}

	return session, nil
}
