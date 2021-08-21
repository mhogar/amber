package migrations

import (
	"authserver/common"
	"authserver/data"
	sqladapter "authserver/data/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

func m003(exec data.DataExecutor, sf data.ScopeFactory) migrationrunner.Migration {
	return migrationrunner.Migration{
		Timestamp:   "003",
		Description: "create sessions table",
		Migrator: &migrator003{
			Executor:     exec,
			ScopeFactory: sf,
		},
	}
}

type migrator003 struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

func (m migrator003) Up() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//create the session table
		err := sqlTx.CreateSessionTable()
		if err != nil {
			return false, common.ChainError("error creating session table", err)
		}

		return true, nil
	})
}

func (m migrator003) Down() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//drop the session table
		err := sqlTx.DropSessionTable()
		if err != nil {
			return false, common.ChainError("error dropping session table", err)
		}

		return true, nil
	})
}
