package migrations

import (
	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
	sqladapter "github.com/mhogar/amber/data/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

func m004(exec data.DataExecutor, sf data.ScopeFactory) migrationrunner.Migration {
	return migrationrunner.Migration{
		Timestamp:   "004",
		Description: "create user-roles table",
		Migrator: &migrator004{
			Executor:     exec,
			ScopeFactory: sf,
		},
	}
}

type migrator004 struct {
	Executor     data.DataExecutor
	ScopeFactory data.ScopeFactory
}

func (m migrator004) Up() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//create the user-role table
		err := sqlTx.CreateUserRoleTable()
		if err != nil {
			return false, common.ChainError("error creating user-role table", err)
		}

		return true, nil
	})
}

func (m migrator004) Down() error {
	return m.ScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//drop the user-role table
		err := sqlTx.DropUserRoleTable()
		if err != nil {
			return false, common.ChainError("error dropping user-role table", err)
		}

		return true, nil
	})
}
