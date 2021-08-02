package dependencies

import (
	"authserver/data"
	"sync"
)

var createScopeFactoryOnce sync.Once
var scopeFactory data.ScopeFactory

// ResolveScopeFactory resolves the ScopeFactory dependency
// Only the first call to this function will create a new ScopeFactory, after which it will be retrieved from memory
func ResolveScopeFactory() data.ScopeFactory {
	createScopeFactoryOnce.Do(func() {
		scopeFactory = &data.CoreScopeFactory{
			DataAdapter: ResolveDataAdapter(),
		}
	})
	return scopeFactory
}
