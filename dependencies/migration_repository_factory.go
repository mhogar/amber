package dependencies

import (
	"sync"

	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/data/database/sql_adapter/migrations"
	firestoreadapter "github.com/mhogar/amber/data/firestore_adapter"
)

var createMigrationRepositoryFactoryOnce sync.Once
var migrationRepositoryFactory data.MigrationRepositoryFactory

// ResolveMigrationRepositoryFactory resolves the MigrationRepositoryFactory dependency.
// Only the first call to this function will create a new MigrationRepositoryFactory, after which it will be retrieved from memory.
func ResolveMigrationRepositoryFactory() data.MigrationRepositoryFactory {
	createMigrationRepositoryFactoryOnce.Do(func() {
		switch config.GetDataAdapter() {
		case "database":
			migrationRepositoryFactory = &migrations.SQLMigrationRepositoryFactory{
				ScopeFactory: ResolveScopeFactory(),
			}
		case "firestore":
			migrationRepositoryFactory = &firestoreadapter.FirestoreMigrationRepositoryFactory{}
		default:
			panic("invalid data adpater key")
		}
	})
	return migrationRepositoryFactory
}
