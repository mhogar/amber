package data

import (
	"github.com/mhogar/amber/common"
)

type ScopeFactory interface {
	// CreateDataExecutorScope handles setup and clean up of the data executor, and passes it to the provided body function.
	// Returns any errors from the executor or body.
	CreateDataExecutorScope(func(DataExecutor) error) error

	// CreateTransactionScope handles transaction creation and commiting, and passes it to the provided body function.
	// Rollbacks on failure or error from body, and returns any errors from the transaction or body.
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
