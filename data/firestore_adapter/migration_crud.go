package firestoreadapter

import (
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"google.golang.org/api/iterator"
)

func (crud *FirestoreCRUD) Setup() error {
	return nil
}

func (crud *FirestoreCRUD) CreateMigration(timestamp string) error {
	//create and validate migration model
	migration := models.CreateMigration(timestamp)
	verr := migration.Validate()
	if verr != models.ValidateMigrationValid {
		return errors.New(fmt.Sprint("error validating migration model:", verr))
	}

	//create the migration
	err := crud.DocWriter.Create(crud.getMigrationDocRef(timestamp), migration)
	if err != nil {
		return common.ChainError("error creating migration", err)
	}

	return nil
}

func (crud *FirestoreCRUD) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	doc, err := crud.getMigration(timestamp)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	return crud.readMigrationData(doc)
}

func (crud *FirestoreCRUD) GetLatestTimestamp() (string, bool, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	itr := crud.Client.Collection("migrations").
		OrderBy("timestamp", firestore.Desc).
		Limit(1).
		Documents(ctx)

	defer cancel()
	defer itr.Stop()

	//get the result
	doc, err := itr.Next()
	if err == iterator.Done {
		return "", false, nil
	}
	if err != nil {
		return "", false, common.ChainError("error getting next doc", err)
	}

	//read the data
	migration, err := crud.readMigrationData(doc)
	if err != nil {
		return "", false, err
	}

	return migration.Timestamp, true, nil
}

func (crud *FirestoreCRUD) DeleteMigrationByTimestamp(timestamp string) error {
	//check migration already exists
	doc, err := crud.getMigration(timestamp)
	if err != nil {
		return err
	}
	if doc == nil {
		return nil
	}

	//delete migration
	err = crud.DocWriter.Delete(doc.Ref)
	if err != nil {
		return common.ChainError("error deleting migration", err)
	}

	return nil
}

func (crud *FirestoreCRUD) getMigrationDocRef(timestamp string) *firestore.DocumentRef {
	return crud.Client.Collection("migrations").Doc(timestamp)
}

func (crud *FirestoreCRUD) getMigration(timestamp string) (*firestore.DocumentSnapshot, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	doc, err := crud.getMigrationDocRef(timestamp).Get(ctx)
	cancel()

	//check migration was found
	if !doc.Exists() {
		return nil, nil
	}

	//handle other errors
	if err != nil {
		return nil, common.ChainError("error getting migration", err)
	}

	return doc, nil
}

func (*FirestoreCRUD) readMigrationData(doc *firestore.DocumentSnapshot) (*models.Migration, error) {
	migration := &models.Migration{}

	err := doc.DataTo(&migration)
	if err != nil {
		return nil, common.ChainError("error reading migration data", err)
	}

	return migration, nil
}
