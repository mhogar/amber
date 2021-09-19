package sqladapter

import (
	"authserver/common"
	"authserver/models"
	"errors"
	"fmt"
)

// Setup creates the migration table if it does not already exist.
func (crud *SQLCRUD) Setup() error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.CreateMigrationTableScript())
	cancel()

	if err != nil {
		return common.ChainError("error executing create migration table statment", err)
	}

	return nil
}

func (crud *SQLCRUD) CreateMigration(timestamp string) error {
	//create and validate migration model
	migration := models.CreateMigration(timestamp)
	verr := migration.Validate()
	if verr != models.ValidateMigrationValid {
		return errors.New(fmt.Sprint("error validating migration model:", verr))
	}

	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.SaveMigrationScript(), migration.Timestamp)
	cancel()

	if err != nil {
		return common.ChainError("error executing save migration statment", err)
	}

	return nil
}

func (crud *SQLCRUD) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetMigrationByTimestampScript(), timestamp)
	defer cancel()

	if err != nil {
		return nil, common.ChainError("error executing get migration by timestamp query", err)
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, common.ChainError("error preparing next row", err)
		}

		//return no results
		return nil, nil
	}

	//get the result
	migration := &models.Migration{}
	err = rows.Scan(&migration.Timestamp)
	if err != nil {
		return nil, common.ChainError("error reading row", err)
	}

	return migration, nil
}

func (crud *SQLCRUD) GetLatestTimestamp() (timestamp string, hasLatest bool, err error) {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	rows, err := crud.Executor.QueryContext(ctx, crud.SQLDriver.GetLatestTimestampScript())
	defer cancel()

	if err != nil {
		return "", false, common.ChainError("error executing get latest timestamp query", err)
	}
	defer rows.Close()

	//check if there was a result
	if !rows.Next() {
		err := rows.Err()
		if err != nil {
			return "", false, common.ChainError("error preparing next row", err)
		}

		//return no results
		return "", false, nil
	}

	//get the result
	err = rows.Scan(&timestamp)
	if err != nil {
		return "", false, common.ChainError("error reading row", err)
	}

	return timestamp, true, nil
}

func (crud *SQLCRUD) DeleteMigrationByTimestamp(timestamp string) error {
	ctx, cancel := crud.ContextFactory.CreateStandardTimeoutContext()
	_, err := crud.Executor.ExecContext(ctx, crud.SQLDriver.DeleteMigrationByTimestampScript(), timestamp)
	cancel()

	if err != nil {
		return common.ChainError("error executing delete migration by timestamp statement", err)
	}

	return nil
}
