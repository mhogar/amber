package sqladapter

import (
	"context"
	"database/sql"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/data"
)

type contextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type SQLCRUD struct {
	Executor       contextExecutor
	SQLDriver      SQLDriver
	ContextFactory data.ContextFactory
}

type SQLTransaction struct {
	*sql.Tx
	SQLCRUD
}

type SQLExecutor struct {
	DB *sql.DB
	SQLCRUD
}

func (exec *SQLExecutor) CreateTransaction() (data.Transaction, error) {
	tx, err := exec.DB.Begin()
	if err != nil {
		return nil, common.ChainError("error beginning transaction", err)
	}

	return &SQLTransaction{
		Tx: tx,
		SQLCRUD: SQLCRUD{
			Executor:       tx,
			SQLDriver:      exec.SQLDriver,
			ContextFactory: exec.ContextFactory,
		},
	}, nil
}
