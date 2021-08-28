package loaders

import (
	"authserver/common"
	"authserver/config"
	"crypto/rsa"
	"io/ioutil"
	"path"

	"github.com/golang-jwt/jwt"
)

type StaticRSAKeyLoader struct{}

// LoadPrivateKey loads the private key from the uri in the static directory in the project.
// Returns any errors.
func (StaticRSAKeyLoader) LoadPrivateKey(uri string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(path.Join(config.GetAppRoot(), "static", uri))
	if err != nil {
		return nil, common.ChainError("error reading key file", err)
	}

	return jwt.ParseRSAPrivateKeyFromPEM(bytes)
}
