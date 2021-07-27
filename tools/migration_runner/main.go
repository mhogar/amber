package main

import (
	"authserver/common"
	"authserver/config"
	"authserver/data"
	"authserver/dependencies"
	"authserver/tools/migration_runner/interfaces"
	"flag"
	"log"

	"github.com/spf13/viper"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	//parse flags
	down := flag.Bool("down", false, "Run migrate down instead of migrate up")
	dbKey := flag.String("db", "core", "The database to run the migrations against")
	flag.Parse()

	viper.Set("db_key", *dbKey)

	mrf := interfaces.MigrationRunnerFactory{
		MigrationRepositoryFactory: dependencies.ResolveMigrationRepositoryFactory(),
	}

	err = Run(dependencies.ResolveScopeFactory(), &mrf, *down)
	if err != nil {
		log.Fatal(err)
	}
}

// Run runs the migration runner. Returns any errors.
func Run(sf data.IScopeFactory, mrf interfaces.IMigrationRunnerFactory, down bool) error {
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
