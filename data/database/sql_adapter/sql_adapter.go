package sqladapter

import (
	"context"
	"database/sql"

	"github.com/mhogar/amber/data"
	"github.com/mhogar/amber/data/database"
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
	ContextFactory data.ContextFactory

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

// CreateSQLAdpater creates a new SQLAdapter with the provided db key and driver.
func CreateSQLAdpater(dbKey string, driver SQLDriver) *SQLAdapter {
	adapter := &SQLAdapter{
		DbKey:     dbKey,
		SQLDriver: driver,
	}
	adapter.Connection = adapter

	return adapter
}
