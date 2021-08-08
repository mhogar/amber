package dependencies

import (
	"authserver/router"
	"sync"
)

var createRouterFactoryOnce sync.Once
var routerFactory router.CoreRouterFactory

// ResolveRouterFactory resolves the CoreRouterFactory dependency
// Only the first call to this function will create a new CoreRouterFactory, after which it will be retrieved from memory
func ResolveRouterFactory() router.RouterFactory {
	createRouterFactoryOnce.Do(func() {
		routerFactory = router.CoreRouterFactory{
			CoreScopeFactory: ResolveScopeFactory(),
			Handlers:         ResolveHandlers(),
		}
	})
	return routerFactory
}
