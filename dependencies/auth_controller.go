package dependencies

import (
	controllerspkg "authserver/controllers"
	"sync"
)

var createAuthControllerOnce sync.Once
var authController controllerspkg.AuthController

// ResolveControllers resolves the AuthController dependency.
// Only the first call to this function will create a new AuthController, after which it will be retrieved from memory.
func ResolveAuthController() controllerspkg.AuthController {
	createAuthControllerOnce.Do(func() {
		authController = &controllerspkg.CoreAuthController{
			PasswordHasher: ResolvePasswordHasher(),
		}
	})
	return authController
}
