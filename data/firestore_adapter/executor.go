package firestoreadapter

import (
	"context"

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

func (*FirestoreExecutor) Create(ctx context.Context, ref *firestore.DocumentRef, data interface{}) error {
	_, err := ref.Create(ctx, data)
	return err
}

func (*FirestoreExecutor) Update(ctx context.Context, ref *firestore.DocumentRef, updates []firestore.Update) error {
	_, err := ref.Update(ctx, updates)
	return err
}

func (*FirestoreExecutor) Delete(ctx context.Context, ref *firestore.DocumentRef) error {
	_, err := ref.Delete(ctx)
	return err
}
