package main

import (
	"authserver/data"

	"github.com/mhogar/migrationrunner"
)

// MigrationRunner is an interface to match the signature of migrationrunner's MigrationRunner.
type MigrationRunner interface {
	MigrateUp() error
	MigrateDown() error
}

type IMigrationRunnerFactory interface {
	CreateMigrationRunner(data.DataExecutor) MigrationRunner
}

type MigrationRunnerFactory struct {
	MigrationRepositoryFactory data.MigrationRepositoryFactory
}

func (mrf MigrationRunnerFactory) CreateMigrationRunner(exec data.DataExecutor) MigrationRunner {
	return migrationrunner.MigrationRunner{
		MigrationRepository: mrf.MigrationRepositoryFactory.CreateMigrationRepository(exec),
		MigrationCRUD:       exec,
	}
}
