package database

import "github.com/mhogar/amber/common"

type DBConnection interface {
	// OpenConnection opens the connection to the database. Returns any errors.
	OpenConnection() error

	// CloseConnection closes the connection to the database and cleans up associated resources. Returns any errors.
	CloseConnection() error

	// Ping pings the database to verify it can be reached.
	// Returns an error if the database can't be reached or if any other errors occur.
	Ping() error
}

type DatabaseAdapter struct {
	Connection DBConnection
}

func (db *DatabaseAdapter) Setup() error {
	//open the db connection
	err := db.Connection.OpenConnection()
	if err != nil {
		return common.ChainError("could not open database connection", err)
	}

	//check db is connected
	err = db.Connection.Ping()
	if err != nil {
		return common.ChainError("could not reach database", err)
	}

	return nil
}

func (db *DatabaseAdapter) CleanUp() error {
	return db.Connection.CloseConnection()
}
