package dependencies

import (
	"sync"

	"github.com/mhogar/amber/loaders"
)

var createJSONLoaderOnce sync.Once
var jsonLoader loaders.JSONLoader

// ResolveJSONLoader resolves the JSONLoader dependency.
// Only the first call to this function will create a new JSONLoader, after which it will be retrieved from memory.
func ResolveJSONLoader() loaders.JSONLoader {
	createJSONLoaderOnce.Do(func() {
		jsonLoader = loaders.StaticJSONLoader{}
	})
	return jsonLoader
}
