package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/common"
)

type FirestoreTransaction struct {
	FirestoreExecutor
	Batch *firestore.WriteBatch
}

func (t *FirestoreTransaction) Commit() error {
	ctx, cancel := t.ContextFactory.CreateStandardTimeoutContext()
	_, err := t.Batch.Commit(ctx)
	cancel()

	if err != nil {
		return common.ChainError("error commiting batch", err)
	}

	return nil
}

func (*FirestoreTransaction) Rollback() error {
	return nil
}
