package firestoreadapter

import (
	"github.com/mhogar/amber/data"
	"github.com/mhogar/migrationrunner"
)

type MigrationRepository struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

func (repo MigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{}
}

type FirestoreMigrationRepositoryFactory struct {
	ScopeFactory data.ScopeFactory
}

func (FirestoreMigrationRepositoryFactory) CreateMigrationRepository(_ data.DataExecutor) migrationrunner.MigrationRepository {
	return &MigrationRepository{}
}
