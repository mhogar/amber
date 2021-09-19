package dependencies

import (
	"sync"

	jwthelpers "github.com/mhogar/amber/controllers/jwt_helpers"
)

var createTokenSignerOnce sync.Once
var tokenSigner jwthelpers.TokenSigner

// ResolveTokenSigner resolves the TokenSigner dependency.
// Only the first call to this function will create a new TokenSigner, after which it will be retrieved from memory.
func ResolveTokenSigner() jwthelpers.TokenSigner {
	createTokenSignerOnce.Do(func() {
		tokenSigner = jwthelpers.JWTTokenSigner{}
	})
	return tokenSigner
}
