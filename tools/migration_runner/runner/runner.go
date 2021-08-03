package runner

import (
	"authserver/common"
	"authserver/data"
)

// Run runs the migration runner. Returns any errors.
func Run(sf data.ScopeFactory, mrf MigrationRunnerFactory, down bool) error {
	return sf.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		mr := mrf.CreateMigrationRunner(exec)
		var err error

		//run the migrations
		if down {
			err = mr.MigrateDown()
		} else {
			err = mr.MigrateUp()
		}

		if err != nil {
			return common.ChainError("error running migrations", err)
		}
		return nil
	})
}
