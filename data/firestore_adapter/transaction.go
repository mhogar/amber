package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/common"
)

type FirestoreTransaction struct {
	FirestoreCRUD
	Batch *firestore.WriteBatch

	hasWrites bool
}

func (tx *FirestoreTransaction) Create(ref *firestore.DocumentRef, data interface{}) error {
	tx.hasWrites = true
	tx.Batch.Create(ref, data)
	return nil
}

func (tx *FirestoreTransaction) Set(ref *firestore.DocumentRef, data interface{}) error {
	tx.hasWrites = true
	tx.Batch.Set(ref, data)
	return nil
}

func (tx *FirestoreTransaction) Update(ref *firestore.DocumentRef, updates []firestore.Update) error {
	tx.hasWrites = true
	tx.Batch.Update(ref, updates)
	return nil
}

func (tx *FirestoreTransaction) Delete(ref *firestore.DocumentRef) error {
	tx.hasWrites = true
	tx.Batch.Delete(ref)
	return nil
}

func (tx *FirestoreTransaction) Commit() error {
	if !tx.hasWrites {
		return nil
	}

	ctx, cancel := tx.ContextFactory.CreateStandardTimeoutContext()
	_, err := tx.Batch.Commit(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error commiting batch", err)
	}

	return nil
}

func (*FirestoreTransaction) Rollback() error {
	return nil
}
