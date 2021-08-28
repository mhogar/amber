package dependencies

import (
	"authserver/loaders"
	"sync"
)

var createRSAKeyLoaderOnce sync.Once
var rsaKeyLoader loaders.RSAKeyLoader

// ResolveRSAKeyLoader resolves the RSAKeyLoader dependency.
// Only the first call to this function will create a new RSAKeyLoader, after which it will be retrieved from memory.
func ResolveRSAKeyLoader() loaders.RSAKeyLoader {
	createRSAKeyLoaderOnce.Do(func() {
		rsaKeyLoader = loaders.StaticRSAKeyLoader{}
	})
	return rsaKeyLoader
}
