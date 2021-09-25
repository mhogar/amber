package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mhogar/amber/config"

	"gopkg.in/yaml.v3"
)

func main() {
	name := flag.String("name", "local", "The name of the config file")
	flag.Parse()

	err := Run(*name)
	if err != nil {
		log.Fatal(err)
	}
}

// Run runs the config generator and returns any errors.
func Run(name string) error {
	filename := fmt.Sprint("config.", name, ".yml")

	//check if file already exists
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return errors.New("file already exists")
	}

	//create the config struct
	cfg := config.Config{
		AppName:     "Amber",
		DataAdapter: "database",
		TokenConfig: config.TokenConfig{
			DefaultIssuer: "amber",
			Lifetime:      60,
		},
		PermissionConfig: config.PermissionConfig{
			MinClientRank: 5,
		},
		DatabaseConfig: config.DatabaseConfig{
			ConnectionStrings: map[string]string{
				"core":        "",
				"integration": "",
			},
			Timeout: 5000,
		},
		FirestoreConfig: config.FirestoreConfig{
			ServiceFile: "",
			Timeout:     5000,
		},
		PasswordCriteriaConfig: config.PasswordCriteriaConfig{
			MinLength:        8,
			RequireLowerCase: true,
			RequireUpperCase: true,
			RequireDigit:     true,
			RequireSymbol:    true,
		},
	}

	//marshal into yaml format
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	//write to file (permissions read and write)
	return ioutil.WriteFile(filename, data, 0666)
}
