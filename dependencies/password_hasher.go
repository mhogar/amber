package dependencies

import (
	"sync"

	passwordhelpers "github.com/mhogar/amber/controllers/password_helpers"
)

var createPasswordHasherOnce sync.Once
var passwordHasher passwordhelpers.PasswordHasher

// ResolvePasswordHasher resolves the PasswordHasher dependency.
// Only the first call to this function will create a new PasswordHasher, after which it will be retrieved from memory.
func ResolvePasswordHasher() passwordhelpers.PasswordHasher {
	createPasswordHasherOnce.Do(func() {
		passwordHasher = passwordhelpers.BCryptPasswordHasher{}
	})
	return passwordHasher
}
