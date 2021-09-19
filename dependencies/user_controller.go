package dependencies

import (
	"sync"

	controllerspkg "github.com/mhogar/amber/controllers"
)

var createUserControllerOnce sync.Once
var userController *controllerspkg.CoreUserController

// ResolveUserController resolves the UserController dependency.
// Only the first call to this function will create a new UserController, after which it will be retrieved from memory.
func ResolveUserController() controllerspkg.UserController {
	createUserControllerOnce.Do(func() {
		userController = &controllerspkg.CoreUserController{
			PasswordHasher:            ResolvePasswordHasher(),
			PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			AuthController:            ResolveAuthController(),
		}
		userController.UserController = userController
	})
	return userController
}
