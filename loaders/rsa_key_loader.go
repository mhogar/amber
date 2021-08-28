package loaders

import (
	"authserver/common"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type RSAKeyLoader interface {
	// LoadPrivateKeyFromBytes loads the RSA private key from the provided byte splice.
	// Returns the key and any errors.
	LoadPrivateKeyFromBytes(key []byte) (*rsa.PrivateKey, error)

	// LoadPrivateKeyFromString loads the RSA private key from the provided uri.
	// Returns the key and any errors.
	LoadPrivateKeyFromURI(url string) (*rsa.PrivateKey, error)
}

type RSAKeyLoaderBase struct{}

func (RSAKeyLoaderBase) LoadPrivateKeyFromBytes(key []byte) (*rsa.PrivateKey, error) {
	//parse the key string
	block, _ := pem.Decode(key)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, common.ChainError("error parsing key string string", err)
	}

	return privateKey.(*rsa.PrivateKey), nil
}
