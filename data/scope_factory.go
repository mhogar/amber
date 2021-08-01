package data

import (
	"authserver/common"
)

type ScopeFactory interface {
	CreateDataExecutorScope(func(DataExecutor) error) error
	CreateTransactionScope(DataExecutor, func(Transaction) (bool, error)) error
}

type CoreScopeFactory struct {
	DataAdapter DataAdapter
}

func (sf CoreScopeFactory) CreateDataExecutorScope(body func(DataExecutor) error) error {
	//init the data adapter
	err := sf.DataAdapter.Setup()
	if err != nil {
		return common.ChainError("error setting up data adapter", err)
	}
	defer sf.DataAdapter.CleanUp()

	//execute the body
	return body(sf.DataAdapter.GetExecutor())
}

func (sf CoreScopeFactory) CreateTransactionScope(exec DataExecutor, body func(Transaction) (bool, error)) error {
	tx, err := exec.CreateTransaction()
	if err != nil {
		return common.ChainError("error creating transaction", err)
	}
	defer tx.Rollback()

	//execute the body
	result, err := body(tx)
	if err != nil {
		return err
	}

	//commit the transaction if success
	if result {
		err = tx.Commit()
		if err != nil {
			return common.ChainError("error commiting transaction", err)
		}
	}

	return nil
}
