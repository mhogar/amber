package runner

import (
	"authserver/data"

	"github.com/mhogar/migrationrunner"
)

// MigrationRunner is an interface to match the signature of migrationrunner's MigrationRunner
type MigrationRunner interface {
	MigrateUp() error
	MigrateDown() error
}

type MigrationRunnerFactory interface {
	// CreateMigrationRunner creates a MigrationRunner using the provided data executor
	CreateMigrationRunner(data.DataExecutor) MigrationRunner
}

type CoreMigrationRunnerFactory struct {
	MigrationRepositoryFactory data.MigrationRepositoryFactory
}

func (mrf CoreMigrationRunnerFactory) CreateMigrationRunner(exec data.DataExecutor) MigrationRunner {
	return migrationrunner.MigrationRunner{
		MigrationRepository: mrf.MigrationRepositoryFactory.CreateMigrationRepository(exec),
		MigrationCRUD:       exec,
	}
}
