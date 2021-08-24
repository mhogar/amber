package loaders

import (
	"authserver/common"
	"authserver/config"
	"encoding/json"
	"os"
	"path"
)

type StaticJSONLoader struct{}

//Load loads the json from the static directory in the project.
//Returns any errors.
func (StaticJSONLoader) Load(uri string, v interface{}) error {
	//open the json file
	file, err := os.Open(path.Join(config.GetAppRoot(), "static", uri))
	if err != nil {
		return common.ChainError("error opening json file", err)
	}
	defer file.Close()

	// decode the json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(v)
	if err != nil {
		return common.ChainError("invalid json file", err)
	}

	return nil
}
