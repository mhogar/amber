package migrations

import (
	"authserver/data"

	"github.com/mhogar/migrationrunner"
)

// MigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the sql db.
type MigrationRepository struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo MigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		m001(repo.Executor, repo.ScopeFactory),
		m002(repo.Executor, repo.ScopeFactory),
		m003(repo.Executor, repo.ScopeFactory),
		m004(repo.Executor, repo.ScopeFactory),
	}
}

type SQLMigrationRepositoryFactory struct {
	ScopeFactory data.ScopeFactory
}

func (f SQLMigrationRepositoryFactory) CreateMigrationRepository(exec data.DataExecutor) migrationrunner.MigrationRepository {
	return &MigrationRepository{
		Executor:     exec,
		ScopeFactory: f.ScopeFactory,
	}
}
