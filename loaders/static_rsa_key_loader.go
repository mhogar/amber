package loaders

import (
	"authserver/common"
	"authserver/config"
	"crypto/rsa"
	"io/ioutil"
	"path"
)

type StaticRSAKeyLoader struct {
	RSAKeyLoaderBase
}

// LoadPrivateKeyFromURI loads the private key from the uri in the static directory in the project.
// Returns any errors.
func (l StaticRSAKeyLoader) LoadPrivateKeyFromURI(uri string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(path.Join(config.GetAppRoot(), "static", uri))
	if err != nil {
		return nil, common.ChainError("error reading key file", err)
	}

	return l.LoadPrivateKeyFromBytes(bytes)
}
