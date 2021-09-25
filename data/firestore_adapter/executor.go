package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/data"
)

type FirestoreExecutor struct {
	FirestoreCRUD
	Client *firestore.Client
}

func (exec *FirestoreExecutor) CreateTransaction() (data.Transaction, error) {
	return &FirestoreTransaction{
		FirestoreExecutor: *exec,
		Batch:             exec.Client.Batch(),
	}, nil
}
