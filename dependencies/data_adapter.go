package dependencies

import (
	"sync"

	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/data"
	sqladapter "github.com/mhogar/amber/data/database/sql_adapter"
	firestoreadapter "github.com/mhogar/amber/data/firestore_adapter"

	"github.com/spf13/viper"
)

var createDataApdaterOnce sync.Once
var dataAdapter data.DataAdapter

// ResolveDataAdapter resolves the DataAdapter dependency.
// Only the first call to this function will create a new DataAdapter, after which it will be retrieved from memory.
func ResolveDataAdapter() data.DataAdapter {
	createDataApdaterOnce.Do(func() {
		switch config.GetDataAdapter() {
		case "database":
			dataAdapter = sqladapter.CreateSQLAdpater(viper.GetString("db_key"), ResolveSQLDriver())
		case "firestore":
			dataAdapter = &firestoreadapter.FirestoreAdapter{}
		default:
			panic("invalid data adapter key")
		}
	})
	return dataAdapter
}
