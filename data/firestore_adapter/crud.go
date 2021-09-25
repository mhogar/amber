package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/data"
)

type FirestoreCRUD struct {
	*firestore.Client
	ContextFactory data.ContextFactory
}
