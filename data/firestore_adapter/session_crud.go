package firestoreadapter

import (
	"errors"
	"fmt"

	"github.com/mhogar/amber/models"

	"github.com/google/uuid"
)

func (crud *FirestoreCRUD) SaveSession(session *models.Session) error {
	//validate the session model
	verr := session.Validate()
	if verr != models.ValidateSessionValid {
		return errors.New(fmt.Sprint("error validating session model:", verr))
	}

	return nil
}

func (crud *FirestoreCRUD) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) DeleteSession(token uuid.UUID) (bool, error) {
	return false, nil
}

func (crud *FirestoreCRUD) DeleteAllUserSessions(username string) error {
	return nil
}

func (crud *FirestoreCRUD) DeleteAllOtherUserSessions(username string, token uuid.UUID) error {
	return nil
}
