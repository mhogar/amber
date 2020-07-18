package sqladapter

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"errors"
	"fmt"
)

// Setup creates the migration table if it does not already exist.
func (adapter *SQLAdapter) Setup() error {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.CreateMigrationTableScript())
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing create migration table statment", err)
	}

	return nil
}

// CreateMigration validates the given timestamp and inserts it into the migration table.
// Returns any errors.
func (adapter *SQLAdapter) CreateMigration(timestamp string) error {
	//create and validate migration model
	migration := models.CreateNewMigration(timestamp)
	verr := migration.Validate()
	if verr != models.ValidateMigrationValid {
		return errors.New(fmt.Sprint("error validating migration model:", verr))
	}

	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.SaveMigrationScript(), migration.Timestamp)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing save migration statment", err)
	}

	return nil
}

// GetMigrationByTimestamp gets the row in the migration table with the matching timestamp, and creates a new migration model using its data.
// Returns the model and any errors.
func (adapter *SQLAdapter) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetMigrationByTimestampScript(), timestamp)
	defer cancel()

	if err != nil {
		return nil, commonhelpers.ChainError("error executing get migration by timestamp query", err)
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, commonhelpers.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	//get the result
	migration := &models.Migration{}
	err = rows.Scan(&migration.Timestamp)
	if err != nil {
		return nil, commonhelpers.ChainError("error reading row", err)
	}

	return migration, nil
}

// GetLatestTimestamp returns the latest timestamp of all rows in the migration table.
// If the table is empty, hasLatest will be false, else it will be true.
// Returns any errors.
func (adapter *SQLAdapter) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	rows, err := adapter.SQLExecuter.QueryContext(ctx, adapter.SQLDriver.GetLatestTimestampScript())
	defer cancel()

	if err != nil {
		return "", false, commonhelpers.ChainError("error executing get latest timestamp query", err)
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return "", false, commonhelpers.ChainError("error preparing next row", err)
		}

		//return no results
		return "", false, nil
	}

	//get the result
	err = rows.Scan(&timestamp)
	if err != nil {
		return "", false, commonhelpers.ChainError("error reading row", err)
	}

	return timestamp, true, nil
}

// DeleteMigrationByTimestamp deletes up to one row from the migartion table with a matching timestamp.
// Returns any errors.
func (adapter *SQLAdapter) DeleteMigrationByTimestamp(timestamp string) error {
	ctx, cancel := adapter.CreateStandardTimeoutContext()
	_, err := adapter.SQLExecuter.ExecContext(ctx, adapter.SQLDriver.DeleteMigrationByTimestampScript(), timestamp)
	cancel()

	if err != nil {
		return commonhelpers.ChainError("error executing delete migration by timestamp statement", err)
	}

	return nil
}
