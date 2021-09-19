package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/mhogar/amber/data/database/sql_adapter/postgres"
)

func main() {
	//parse flags
	name := flag.String("name", "", "The name of the database to create")
	user := flag.String("U", "postgres", "The database user")
	password := flag.String("P", "password", "The user's password")
	host := flag.String("H", "localhost", "The server's host")
	port := flag.String("p", "5432", "The server's port on the host")
	flag.Parse()

	err := createDatabase(*name, *user, *password, *host, *port)
	if err != nil {
		log.Fatal(err)
	}
}

func createDatabase(name string, user string, password string, host string, port string) error {
	connection_str := fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=disable", user, password, host, port)

	//connect to the db
	db, err := sql.Open(postgres.Driver{}.GetDriverName(), connection_str)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Dropping", name, "if it exists")
	_, err = db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}

	log.Println("Creating database", name)
	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}

	return nil
}
