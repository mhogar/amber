package main

import (
	"authserver/config"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

type tmplData struct {
	Name   string
	Script string
}

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	inDir := path.Join(viper.GetString("root_dir"), "database/sql_adapter/postgres/scripts")

	//load the template
	tmpl := template.Must(template.ParseFiles(path.Join(viper.GetString("root_dir"), "database/sql_adapter/script_compiler/script_repository.go.tmpl")))

	var data []tmplData
	createDataObjects(inDir, &data)

	//create the output file
	f, err := os.Create(path.Join(inDir, "script_repository.go"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//execute the template
	tmpl.Execute(f, data)
}

func createDataObjects(dir string, data *[]tmplData) {
	//get all files in input dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//read all sql files and create data objects
	for _, file := range files {
		//if another directory, recurse
		if file.IsDir() {
			createDataObjects(path.Join(dir, file.Name()), data)
			continue
		}

		//if sql file, create data object and add to slice
		if path.Ext(file.Name()) == ".sql" {
			script, err := ioutil.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			*data = append(*data, tmplData{
				Name:   strings.SplitN(path.Base(file.Name()), ".", 2)[0],
				Script: string(script),
			})
		}
	}
}
