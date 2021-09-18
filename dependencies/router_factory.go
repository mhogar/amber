package dependencies

import (
	"authserver/router"
	"sync"
)

var createRouterFactoryOnce sync.Once
var routerFactory router.RouterFactory

// ResolveRouterFactory resolves the RouterFactory dependency.
// Only the first call to this function will create a new RouterFactory, after which it will be retrieved from memory.
func ResolveRouterFactory() router.RouterFactory {
	createRouterFactoryOnce.Do(func() {
		routerFactory = router.CoreRouterFactory{
			ScopeFactory: ResolveScopeFactory(),
			Handlers:     ResolveHandlers(),
		}
	})
	return routerFactory
}
