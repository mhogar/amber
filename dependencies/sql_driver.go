package dependencies

import (
	"sync"

	"github.com/mhogar/amber/config"
	sqladapter "github.com/mhogar/amber/data/database/sql_adapter"
	"github.com/mhogar/amber/data/database/sql_adapter/postgres"
)

var createSQLDriverOnce sync.Once
var sqlDriver sqladapter.SQLDriver

// ResolveSQLDriver resolves the SQLDriver dependency.
// Only the first call to this function will create a new SQLDriver, after which it will be retrieved from memory.
func ResolveSQLDriver() sqladapter.SQLDriver {
	createSQLDriverOnce.Do(func() {
		switch config.GetDatabaseConfig().Driver {
		case "postgres":
			sqlDriver = postgres.Driver{}
		default:
			panic("invalid database driver")
		}
	})
	return sqlDriver
}
