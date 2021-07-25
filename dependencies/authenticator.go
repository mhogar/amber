package dependencies

import (
	"authserver/router"
	"sync"
)

var createAuthenticatorOnce sync.Once
var authenticator router.Authenticator

// ResolveAuthenticator resolves the Authenticator dependency.
// Only the first call to this function will create a new Authenticator, after which it will be retrieved from memory.
func ResolveAuthenticator() router.Authenticator {
	createAuthenticatorOnce.Do(func() {
		authenticator = &router.OAuthAuthenticator{}
	})
	return authenticator
}
