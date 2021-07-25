package postgres

import (
	"authserver/data/database/sql_adapter/postgres/scripts"

	//import the postgres driver
	_ "github.com/lib/pq"
)

// Driver is an implementation of the SQL Driver interface for postgres.
type Driver struct {
	scripts.ScriptRepository
}

// GetDriverName returns the postgres driver name.
func (Driver) GetDriverName() string {
	return "postgres"
}
