package main

import (
	"authserver/config"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
		AppName: "Amber",
		TokenConfig: config.TokenConfig{
			DefaultIssuer: "amber",
			Lifetime:      60,
		},
		DatabaseConfig: config.DatabaseConfig{
			ConnectionStrings: map[string]string{
				"core":        "",
				"integration": "",
			},
			Timeout: 3000,
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

	//write to file (permissions read & write)
	return ioutil.WriteFile(filename, data, 0666)
}
