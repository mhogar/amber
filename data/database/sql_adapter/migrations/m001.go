package migrations

import (
	"authserver/common"
	"authserver/data"
	sqladapter "authserver/data/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

func m001(exec data.DataExecutor, sf data.ScopeFactory) migrationrunner.Migration {
	return migrationrunner.Migration{
		Timestamp:   "001",
		Description: "create users table",
		Migrator: &migrator001{
			Executor:     exec,
			ScopeFactory: sf,
		},
	}
}

type migrator001 struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

func (m migrator001) Up() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//create the user table
		err := sqlTx.CreateUserTable()
		if err != nil {
			return false, common.ChainError("error creating user table", err)
		}

		return true, nil
	})
}

func (m migrator001) Down() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//drop the user table
		err := sqlTx.DropUserTable()
		if err != nil {
			return false, common.ChainError("error dropping user table", err)
		}

		return true, nil
	})
}
