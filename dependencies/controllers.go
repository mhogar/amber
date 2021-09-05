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
			UserController: controllerspkg.CoreUserController{
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
				AuthController:            ResolveAuthController(),
			},
			ClientController: controllerspkg.CoreClientController{},
			AuthController:   ResolveAuthController(),
			SessionController: controllerspkg.CoreSessionController{
				AuthController: ResolveAuthController(),
			},
			TokenController: controllerspkg.CoreTokenController{
				AuthController:       ResolveAuthController(),
				TokenFactorySelector: ResolveTokenFactorySelector(),
			},
			UserRoleController: controllerspkg.CoreUserRoleController{},
		}
	})
	return controllers
}
