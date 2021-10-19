package loaders

import (
	"encoding/json"
	"os"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
)

type StaticJSONLoader struct{}

// Load loads the json from the static directory in the project.
// Returns any errors.
func (StaticJSONLoader) Load(uri string, v interface{}) error {
	//open the json file
	file, err := os.Open(config.GetAppRoot("static", uri))
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
