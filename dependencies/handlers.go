package dependencies

import (
	handlerspkg "authserver/router/handlers"
	"sync"
)

var createHandlersOnce sync.Once
var handlers handlerspkg.Handlers

// ResolveHandlers resolves the Handlers dependency.
// Only the first call to this function will create a new Handlers, after which it will be retrieved from memory.
func ResolveHandlers() handlerspkg.Handlers {
	createHandlersOnce.Do(func() {
		handlers = handlerspkg.CoreHandlers{
			Controllers: ResolveControllers(),
			Renderer:    ResolveRenderer(),
		}
	})
	return handlers
}
