package dependencies

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"sync"
)

var createTokenFactoryOnce sync.Once
var tokenFactory jwthelpers.TokenFactory

// ResolveTokenFactory resolves the TokenFactory dependency.
// Only the first call to this function will create a new TokenFactory, after which it will be retrieved from memory.
func ResolveTokenFactory() jwthelpers.TokenFactory {
	createTokenFactoryOnce.Do(func() {
		tokenFactory = jwthelpers.FirebaseTokenFactory{}
	})
	return tokenFactory
}
