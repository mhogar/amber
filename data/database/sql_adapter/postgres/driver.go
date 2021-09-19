package postgres

import (
	"github.com/mhogar/amber/data/database/sql_adapter/postgres/scripts"

	//import the postgres driver
	_ "github.com/lib/pq"
)

type Driver struct {
	scripts.ScriptRepository
}

func (Driver) GetDriverName() string {
	return "postgres"
}
