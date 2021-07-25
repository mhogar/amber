package migrations

import (
	"authserver/data"

	"github.com/mhogar/migrationrunner"
)

// MigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the sql db.
type MigrationRepository struct {
	Executor     data.DataExecutor
	ScopeFactory data.IScopeFactory
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo MigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		m20200628151601(repo),
	}
}

type SQLMigrationRepositoryFactory struct {
	ScopeFactory data.IScopeFactory
}

func (f SQLMigrationRepositoryFactory) CreateMigrationRepository(exec data.DataExecutor) migrationrunner.MigrationRepository {
	return &MigrationRepository{
		Executor:     exec,
		ScopeFactory: f.ScopeFactory,
	}
}
