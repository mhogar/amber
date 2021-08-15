package migrations

import (
	"authserver/common"
	"authserver/config"
	"authserver/data"
	sqladapter "authserver/data/database/sql_adapter"
	"authserver/models"
)

type m20200628151601 struct {
	Executor         data.DataExecutor
	CoreScopeFactory data.ScopeFactory
}

func (m m20200628151601) GetTimestamp() string {
	return "20200628151601"
}

func (m m20200628151601) Up() error {
	return m.CoreScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//create the user table
		err := sqlTx.CreateUserTable()
		if err != nil {
			return false, common.ChainError("error creating user table", err)
		}

		//create the client table
		err = sqlTx.CreateClientTable()
		if err != nil {
			return false, common.ChainError("error creating client table", err)
		}

		//add this app as a client
		err = sqlTx.CreateClient(models.CreateClient(config.GetAppId(), "AuthServer"))
		if err != nil {
			return false, common.ChainError("error saving app client", err)
		}

		//create the access_token table
		err = sqlTx.CreateAccessTokenTable()
		if err != nil {
			return false, common.ChainError("error creating scope table", err)
		}

		return true, nil
	})
}

func (m m20200628151601) Down() error {
	return m.CoreScopeFactory.CreateTransactionScope(m.Executor, func(tx data.Transaction) (bool, error) {
		sqlTx := tx.(*sqladapter.SQLTransaction)

		//drop the access token table
		err := sqlTx.DropAccessTokenTable()
		if err != nil {
			return false, common.ChainError("error dropping access token table", err)
		}

		//drop the client table
		err = sqlTx.DropClientTable()
		if err != nil {
			return false, common.ChainError("error dropping client table", err)
		}

		//drop the user table
		err = sqlTx.DropUserTable()
		if err != nil {
			return false, common.ChainError("error dropping user table", err)
		}

		return true, nil
	})
}
