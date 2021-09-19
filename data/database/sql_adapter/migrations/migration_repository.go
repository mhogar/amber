package migrations

import (
	"github.com/mhogar/amber/data"

	"github.com/mhogar/migrationrunner"
)

type MigrationRepository struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

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
