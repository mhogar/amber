package dependencies

import (
	controllerspkg "authserver/controllers"
	"sync"
)

var createContollersOnce sync.Once
var controllers controllerspkg.Controllers

// ResolveControllers resolves the Controllers dependency.
// Only the first call to this function will create a new Controllers, after which it will be retrieved from memory.
func ResolveControllers() controllerspkg.Controllers {
	createContollersOnce.Do(func() {
		controllers = &controllerspkg.CoreControllers{
			CoreUserController: controllerspkg.CoreUserController{
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			},
			CoreClientController: controllerspkg.CoreClientController{},
			CoreTokenController: controllerspkg.CoreTokenController{
				PasswordHasher: ResolvePasswordHasher(),
			},
		}
	})
	return controllers
}
