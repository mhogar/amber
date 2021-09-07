package main

import (
	"authserver/config"
	"authserver/dependencies"
	"authserver/tools/admin_creator/runner"
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
	dbKey := flag.String("db", "core", "The database to run the scipt against")
	username := flag.String("username", "", "The username for the admin")
	password := flag.String("password", "", "The password for the admin")
	rank := flag.Int("rank", 10, "The rank for the admin")
	flag.Parse()

	viper.Set("db_key", *dbKey)

	err = runner.Run(dependencies.ResolveScopeFactory(), dependencies.ResolveControllers(), *username, *password, *rank)
	if err != nil {
		log.Fatal(err)
	}
}
