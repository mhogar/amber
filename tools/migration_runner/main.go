package main

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/tools/migration_runner/runner"
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

	mrf := runner.CoreMigrationRunnerFactory{
		MigrationRepositoryFactory: dependencies.ResolveMigrationRepositoryFactory(),
	}

	err = runner.Run(dependencies.ResolveScopeFactory(), &mrf, *down)
	if err != nil {
		log.Fatal(err)
	}
}
