package dependencies

import (
	rendererpkg "authserver/router/renderer"
	"sync"
)

var createRendererOnce sync.Once
var renderer rendererpkg.Renderer

// ResolveRenderer resolves the Renderer dependency.
// Only the first call to this function will create a new Renderer, after which it will be retrieved from memory.
func ResolveRenderer() rendererpkg.Renderer {
	createRendererOnce.Do(func() {
		renderer = rendererpkg.CoreRenderer{}
	})
	return renderer
}