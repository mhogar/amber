package firestoreadapter

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
)

func (crud *FirestoreCRUD) CreateClient(client *models.Client) error {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	err := crud.DocWriter.Create(ctx, crud.Client.Collection("clients").Doc(client.UID.String()), client)
	cancel()

	if err != nil {
		return common.ChainError("error creating client", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetClients() ([]*models.Client, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) GetClientByUID(uid uuid.UUID) (*models.Client, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) UpdateClient(client *models.Client) (bool, error) {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return false, errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	return false, nil
}

func (crud *FirestoreCRUD) DeleteClient(uid uuid.UUID) (bool, error) {
	return false, nil
}
