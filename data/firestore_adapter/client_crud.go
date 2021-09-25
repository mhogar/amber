package firestoreadapter

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"google.golang.org/api/iterator"
)

func (crud *FirestoreCRUD) CreateClient(client *models.Client) error {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	//create client
	err := crud.DocWriter.Create(crud.Client.Collection("clients").Doc(client.UID.String()), client)
	if err != nil {
		return common.ChainError("error creating client", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetClients() ([]*models.Client, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("clients").
		OrderBy("name", firestore.Asc).
		Documents(ctx)

	defer cancel()
	defer itr.Stop()

	//read the results
	clients := []*models.Client{}
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, common.ChainError("error getting next doc", err)
		}

		client, err := crud.readClientData(doc)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func (crud *FirestoreCRUD) GetClientByUID(uid uuid.UUID) (*models.Client, error) {
	doc, err := crud.getClient(uid)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	return crud.readClientData(doc)
}

func (crud *FirestoreCRUD) UpdateClient(client *models.Client) (bool, error) {
	//validate the client model
	verr := client.Validate()
	if verr != models.ValidateClientValid {
		return false, errors.New(fmt.Sprint("error validating client model: ", verr))
	}

	//check client already exists
	doc, err := crud.getClient(client.UID)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//update client
	err = crud.DocWriter.Set(doc.Ref, client)
	if err != nil {
		return true, common.ChainError("error updating client", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) DeleteClient(uid uuid.UUID) (bool, error) {
	//check client already exists
	doc, err := crud.getClient(uid)
	if err != nil {
		return false, err
	}
	if doc == nil {
		return false, nil
	}

	//delete client
	err = crud.DocWriter.Delete(doc.Ref)
	if err != nil {
		return false, common.ChainError("error deleting client", err)
	}

	return true, nil
}

func (crud *FirestoreCRUD) getClient(uid uuid.UUID) (*firestore.DocumentSnapshot, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	doc, err := crud.Client.Collection("clients").Doc(uid.String()).Get(ctx)
	cancel()

	//check client was found
	if !doc.Exists() {
		return nil, nil
	}

	//handle other errors
	if err != nil {
		return nil, common.ChainError("error getting client", err)
	}

	return doc, nil
}

func (*FirestoreCRUD) readClientData(doc *firestore.DocumentSnapshot) (*models.Client, error) {
	client := &models.Client{}

	err := doc.DataTo(&client)
	if err != nil {
		return nil, common.ChainError("error reading client data", err)
	}

	return client, nil
}
