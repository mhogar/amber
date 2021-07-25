package dependencies

import (
	"authserver/data"
	"sync"
)

var createScopeFactoryOnce sync.Once
var scopeFactory data.IScopeFactory

// ResolveScopeFactory resolves the IScopeFactory dependency.
// Only the first call to this function will create a new IScopeFactory, after which it will be retrieved from memory.
func ResolveScopeFactory() data.IScopeFactory {
	createScopeFactoryOnce.Do(func() {
		scopeFactory = &data.ScopeFactory{
			DataAdapter: ResolveDataAdapter(),
		}
	})
	return scopeFactory
}
