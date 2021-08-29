package dependencies

import (
	"authserver/loaders"
	"sync"
)

var createRawDataLoaderOnce sync.Once
var rawDataLoader loaders.RawDataLoader

// ResolveRawDataLoader resolves the RawDataLoader dependency.
// Only the first call to this function will create a new RawDataLoader, after which it will be retrieved from memory.
func ResolveRawDataLoader() loaders.RawDataLoader {
	createRawDataLoaderOnce.Do(func() {
		rawDataLoader = loaders.StaticRawDataLoader{}
	})
	return rawDataLoader
}
