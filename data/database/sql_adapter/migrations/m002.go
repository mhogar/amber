package migrations

import (
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	sqladapter "github.com/mhogar/amber/data/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

func m002(exec data.DataExecutor, sf data.ScopeFactory) migrationrunner.Migration {
	return migrationrunner.Migration{
		Timestamp:   "002",
		Description: "create clients table",
		Migrator: &migrator002{
			Executor:     exec,
			ScopeFactory: sf,
		},
	}
}

type migrator002 struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

func (m migrator002) Up() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//create the client table
		err := sqlTx.CreateClientTable()
		if err != nil {
			return false, common.ChainError("error creating client table", err)
		}

		return true, nil
	})
}

func (m migrator002) Down() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//drop the client table
		err := sqlTx.DropClientTable()
		if err != nil {
			return false, common.ChainError("error dropping client table", err)
		}

		return true, nil
	})
}
