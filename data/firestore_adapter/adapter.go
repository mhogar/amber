package firestoreadapter

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/data"
	"google.golang.org/api/option"
)

type FirestoreAdapter struct {
	cancelFunc context.CancelFunc

	Client         *firestore.Client
	ContextFactory data.ContextFactory
}

// Setup creates a new firestore client using the firestore config.
// Also initializes the adapter's context and cancel function.
// Returns any errors.
func (a *FirestoreAdapter) Setup() error {
	cfg := config.GetFirestoreConfig()

	//setup context factory
	a.ContextFactory.Context, a.cancelFunc = context.WithCancel(context.Background())
	a.ContextFactory.Timeout = cfg.Timeout

	ctx, cancel := a.ContextFactory.CreateStandardTimeoutContext()
	var err error

	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		a.Client, err = firestore.NewClient(ctx, "emulator-project-id")
	} else {
		credsFile := config.GetAppRoot(cfg.ServiceFile)

		//create the firebase app
		app, appErr := firebase.NewApp(a.ContextFactory.Context, nil, option.WithCredentialsFile(credsFile))
		if appErr != nil {
			return common.ChainError("error creating firebase app", appErr)
		}

		a.Client, err = app.Firestore(ctx)
	}
	cancel()

	if err != nil {
		return common.ChainError("error creating firestore client", err)
	}

	return nil
}

// CleanUp closes the firestore client and reset's the adapter's instance.
// The adapter also calls its cancel function to cancel any child requests that may still be running.
// Neither the adapter's client instance or context should be used after calling this function.
// Returns any errors.
func (a *FirestoreAdapter) CleanUp() error {
	err := a.Client.Close()
	if err != nil {
		return common.ChainError("error closing firestore client", err)
	}

	//cancel any remaining requests that may still be running
	a.cancelFunc()

	//clean up resources
	a.Client = nil

	return nil
}

func (a *FirestoreAdapter) GetExecutor() data.DataExecutor {
	exec := &FirestoreExecutor{
		FirestoreCRUD: FirestoreCRUD{
			Client:         a.Client,
			ContextFactory: a.ContextFactory,
		},
	}
	exec.FirestoreCRUD.DocWriter = exec

	return exec
}
