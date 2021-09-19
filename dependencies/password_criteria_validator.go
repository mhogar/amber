package dependencies

import (
	"sync"

	passwordhelpers "github.com/mhogar/amber/controllers/password_helpers"
)

var createPasswordCriteriaValidatorOnce sync.Once
var passwordCriteriaValidator passwordhelpers.PasswordCriteriaValidator

// ResolvePasswordCriteriaValidator resolves the PasswordCriteriaValidator dependency.
// Only the first call to this function will create a new PasswordCriteriaValidator, after which it will be retrieved from memory.
func ResolvePasswordCriteriaValidator() passwordhelpers.PasswordCriteriaValidator {
	createPasswordCriteriaValidatorOnce.Do(func() {
		passwordCriteriaValidator = passwordhelpers.ConfigPasswordCriteriaValidator{}
	})
	return passwordCriteriaValidator
}
