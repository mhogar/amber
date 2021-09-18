package dependencies

import (
	"authserver/data"
	sqladapter "authserver/data/database/sql_adapter"
	"sync"

	"github.com/spf13/viper"
)

var createDataApdaterOnce sync.Once
var dataAdapter data.DataAdapter

// ResolveDataAdapter resolves the DataAdapter dependency.
// Only the first call to this function will create a new DataAdapter, after which it will be retrieved from memory.
func ResolveDataAdapter() data.DataAdapter {
	createDataApdaterOnce.Do(func() {
		dataAdapter = sqladapter.CreateSQLAdpater(viper.GetString("db_key"), ResolveSQLDriver())
	})
	return dataAdapter
}
