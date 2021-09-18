package loaders

import (
	"authserver/common"
	"authserver/config"
	"io/ioutil"
	"path"
)

type StaticRawDataLoader struct{}

// Load loads the data from the static directory in the project.
// Returns any errors.
func (StaticRawDataLoader) Load(uri string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path.Join(config.GetAppRoot(), "static", uri))
	if err != nil {
		return nil, common.ChainError("error reading file", err)
	}

	return bytes, nil
}
