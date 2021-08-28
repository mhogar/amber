package dependencies

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"sync"
)

var createTokenFactorySelectorOnce sync.Once
var tokenFactorySelector jwthelpers.TokenFactorySelector

// ResolveTokenFactorySelector resolves the TokenFactorySelector dependency.
// Only the first call to this function will create a new TokenFactorySelector, after which it will be retrieved from memory.
func ResolveTokenFactorySelector() jwthelpers.TokenFactorySelector {
	createTokenFactorySelectorOnce.Do(func() {
		tokenFactorySelector = jwthelpers.CoreTokenFactorySelector{
			JSONLoader:  ResolveJSONLoader(),
			TokenSigner: ResolveTokenSigner(),
		}
	})
	return tokenFactorySelector
}
