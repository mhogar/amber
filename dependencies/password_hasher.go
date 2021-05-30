package dependencies

import (
	passwordhelpers "authserver/controllers/password_helpers"
	"sync"
)

var passwordHasherOnce sync.Once
var passwordHasher passwordhelpers.PasswordHasher

// ResolvePasswordHasher resolves the PasswordHasher dependency.
// Only the first call to this function will create a new PasswordHasher, after which it will be retrieved from memory.
func ResolvePasswordHasher() passwordhelpers.PasswordHasher {
	passwordHasherOnce.Do(func() {
		passwordHasher = passwordhelpers.BCryptPasswordHasher{}
	})
	return passwordHasher
}
