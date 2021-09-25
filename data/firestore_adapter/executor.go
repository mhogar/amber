package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/data"
)

type FirestoreExecutor struct {
	FirestoreCRUD
}

func (exec *FirestoreExecutor) CreateTransaction() (data.Transaction, error) {
	tx := &FirestoreTransaction{
		FirestoreCRUD: exec.FirestoreCRUD,
		Batch:         exec.Client.Batch(),
	}
	tx.FirestoreCRUD.DocWriter = tx

	return tx, nil
}

func (exec *FirestoreExecutor) Create(ref *firestore.DocumentRef, data interface{}) error {
	ctx, cancel := exec.ContextFactory.CreateStandardTimeoutContext()
	_, err := ref.Create(ctx, data)
	cancel()

	return err
}

func (exec *FirestoreExecutor) Update(ref *firestore.DocumentRef, updates []firestore.Update) error {
	ctx, cancel := exec.ContextFactory.CreateStandardTimeoutContext()
	_, err := ref.Update(ctx, updates)
	cancel()

	return err
}

func (exec *FirestoreExecutor) Delete(ref *firestore.DocumentRef) error {
	ctx, cancel := exec.ContextFactory.CreateStandardTimeoutContext()
	_, err := ref.Delete(ctx)
	cancel()

	return err
}
