package data

import (
	"authserver/models"

	"github.com/mhogar/migrationrunner"
)

type DataCRUD interface {
	models.MigrationCRUD
	models.UserCRUD
	models.ClientCRUD
	models.ScopeCRUD
	models.AccessTokenCRUD
}

type Transaction interface {
	DataCRUD

	// Commit saves the transaction changes and returns any errors.
	Commit() error

	// Rollback aborts the transaction changes and returns any errors.
	Rollback() error
}

type DataExecutor interface {
	DataCRUD

	// CreateTransaction creates a new transaction and returns any errors.
	CreateTransaction() (Transaction, error)
}

type DataAdapter interface {
	// Setup sets up the adapter and returns any errors.
	Setup() error

	// CleanUp cleans up the adapter and returns any errors.
	CleanUp() error

	// GetExecutor gets the DataExecutor for the adapter.
	GetExecutor() DataExecutor
}

type MigrationRepositoryFactory interface {
	CreateMigrationRepository(DataExecutor) migrationrunner.MigrationRepository
}
