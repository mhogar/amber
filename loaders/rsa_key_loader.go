package loaders

import (
	"crypto/rsa"
)

type RSAKeyLoader interface {
	// LoadPrivateKey loads the RSA private key from the provided uri.
	// Returns the key and any errors.
	LoadPrivateKey(url string) (*rsa.PrivateKey, error)
}
