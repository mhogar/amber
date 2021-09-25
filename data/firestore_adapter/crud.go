package firestoreadapter

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/data"
)

type docWriter interface {
	Create(ctx context.Context, ref *firestore.DocumentRef, data interface{}) error
	Update(ctx context.Context, ref *firestore.DocumentRef, updates []firestore.Update) error
	Delete(ctx context.Context, ref *firestore.DocumentRef) error
}

type FirestoreCRUD struct {
	DocWriter   docWriter
	Client         *firestore.Client
	ContextFactory data.ContextFactory
}
