package dependencies

import (
	"authserver/data"
	"authserver/data/database/sql_adapter/migrations"
	"sync"
)

var createMigrationRepositoryFactoryOnce sync.Once
var migrationRepositoryFactory data.MigrationRepositoryFactory

// ResolveMigrationRepositoryFactory resolves the MigrationRepositoryFactory dependency
// Only the first call to this function will create a new MigrationRepositoryFactory, after which it will be retrieved from memory
func ResolveMigrationRepositoryFactory() data.MigrationRepositoryFactory {
	createMigrationRepositoryFactoryOnce.Do(func() {
		migrationRepositoryFactory = &migrations.SQLMigrationRepositoryFactory{
			CoreScopeFactory: ResolveScopeFactory(),
		}
	})
	return migrationRepositoryFactory
}
