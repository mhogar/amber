package main

import (
	"authserver/common"
	"authserver/config"
	"authserver/controllers"
	"authserver/data"
	"authserver/dependencies"
	"flag"
	"fmt"
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
	flag.Parse()

	viper.Set("db_key", *dbKey)

	err = Run(dependencies.ResolveScopeFactory(), dependencies.ResolveControllers(), *username, *password)
	if err != nil {
		log.Fatal(err)
	}
}

// Run connects to the database and runs the admin creator. Returns any errors.
func Run(sf data.IScopeFactory, c controllers.UserController, username string, password string) error {
	return sf.CreateDataExecutorScope(func(exec data.DataExecutor) error {
		return sf.CreateTransactionScope(exec, func(tx data.Transaction) (bool, error) {
			//save the user
			user, rerr := c.CreateUser(tx, username, password)
			if rerr.Type != common.ErrorTypeNone {
				return false, common.ChainError("error creating user", rerr)
			}

			fmt.Println("Created user:", user.ID.String())
			return true, nil
		})
	})
}
