package dependencies

import (
	"authserver/router"
	"sync"
)

var createHandlersOnce sync.Once
var handlers router.IHandlers

// ResolveHandlers resolves the IHandlers dependency.
// Only the first call to this function will create a new IHandlers, after which it will be retrieved from memory.
func ResolveHandlers() router.IHandlers {
	createHandlersOnce.Do(func() {
		handlers = router.Handlers{
			Controllers: ResolveControllers(),
		}
	})
	return handlers
}
