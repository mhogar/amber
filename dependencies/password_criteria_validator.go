package dependencies

import (
	passwordhelpers "authserver/controllers/password_helpers"
	"sync"
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
