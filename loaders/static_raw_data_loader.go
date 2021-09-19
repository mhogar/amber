package loaders

import (
	"io/ioutil"
	"path"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
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
