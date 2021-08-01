package dependencies

import (
	handlerspkg "authserver/router/handlers"
	"sync"
)

var createHandlersOnce sync.Once
var handlers handlerspkg.IHandlers

// ResolveHandlers resolves the IHandlers dependency.
// Only the first call to this function will create a new IHandlers, after which it will be retrieved from memory.
func ResolveHandlers() handlerspkg.IHandlers {
	createHandlersOnce.Do(func() {
		handlers = handlerspkg.Handlers{
			Controllers: ResolveControllers(),
		}
	})
	return handlers
}
