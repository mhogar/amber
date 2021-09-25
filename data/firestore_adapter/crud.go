package firestoreadapter

import (
	"cloud.google.com/go/firestore"
	"github.com/mhogar/amber/data"
)

type docWriter interface {
	Create(ref *firestore.DocumentRef, data interface{}) error
	Set(ref *firestore.DocumentRef, data interface{}) error
	Update(ref *firestore.DocumentRef, updates []firestore.Update) error
	Delete(ref *firestore.DocumentRef) error
}

type FirestoreCRUD struct {
	DocWriter      docWriter
	Client         *firestore.Client
	ContextFactory data.ContextFactory
}
