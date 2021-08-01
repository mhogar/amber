package sqladapter

import (
	"authserver/data"
	"authserver/data/database"
	"context"
	"database/sql"
)

// SQLDriver is an interface for encapsulating methods specific to each sql driver.
type SQLDriver interface {
	SQLScriptRepository

	// GetDriverName returns the name for the driver.
	GetDriverName() string
}

// SQLAdapter contains methods and members common to the sql db and transaction structs.
type SQLAdapter struct {
	database.DatabaseAdapter
	cancelFunc context.CancelFunc

	// DbKey is the key that will be used to resolve the database's connection string.
	DbKey string

	// DB is the sql database instance.
	DB *sql.DB

	// SQLDriver is a dependency for fetching the sql scripts and resolving the driver name.
	SQLDriver SQLDriver

	ContextFactory ContextFactory
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
