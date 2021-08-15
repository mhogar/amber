package sqladapter

import (
	"authserver/data"
	"authserver/data/database"
	"context"
	"database/sql"
)

type SQLDriver interface {
	SQLScriptRepository

	// GetDriverName returns the name for the driver.
	GetDriverName() string
}

type SQLAdapter struct {
	database.DatabaseAdapter
	cancelFunc context.CancelFunc

	DB             *sql.DB
	SQLDriver      SQLDriver
	ContextFactory ContextFactory

	// DbKey is the key that will be used to resolve the database's connection string.
	DbKey string
}

func (a *SQLAdapter) GetExecutor() data.DataExecutor {
	return &SQLExecutor{
		DB: a.DB,
		SQLCRUD: SQLCRUD{
			Executor:       a.DB,
			SQLDriver:      a.SQLDriver,
			ContextFactory: a.ContextFactory,
		},
	}
}

func CreateSQLAdpater(dbKey string, SQLDriver SQLDriver) *SQLAdapter {
	adapter := &SQLAdapter{
		DbKey:     dbKey,
		SQLDriver: SQLDriver,
	}
	adapter.Connection = adapter

	return adapter
}
