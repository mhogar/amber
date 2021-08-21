package dependencies

import (
	controllerspkg "authserver/controllers"
	"sync"
)

var createControllersOnce sync.Once
var controllers controllerspkg.Controllers

// ResolveControllers resolves the Controllers dependency.
// Only the first call to this function will create a new Controllers, after which it will be retrieved from memory.
func ResolveControllers() controllerspkg.Controllers {
	createControllersOnce.Do(func() {
		controllers = &controllerspkg.CoreControllers{
			CoreUserController: controllerspkg.CoreUserController{
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			},
			CoreClientController: controllerspkg.CoreClientController{},
			CoreAuthController:   ResolveAuthController().(controllerspkg.CoreAuthController),
			CoreSessionController: controllerspkg.CoreSessionController{
				AuthController: ResolveAuthController(),
			},
		}
	})
	return controllers
}
