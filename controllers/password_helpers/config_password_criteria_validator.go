package passwordhelpers

import (
	"authserver/config"
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

// ConfigPasswordCriteriaValidator is an implementation of PasswordCriteriaValidator that uses criteria loaded from config
type ConfigPasswordCriteriaValidator struct{}

func (ConfigPasswordCriteriaValidator) ValidatePasswordCriteria(password string) ValidatePasswordCriteriaError {
	criteria := viper.Get("password_criteria").(config.PasswordCriteriaConfig)

	//validate min length
	if len(password) < criteria.MinLength {
		return CreateValidatePasswordCriteriaError(ValidatePasswordCriteriaTooShort, fmt.Sprintf("password must be at least %d characters", criteria.MinLength))
	}

	//validate has lower case letter if required
	if criteria.RequireLowerCase {
		matched, _ := regexp.Match("[a-z]", []byte(password))
		if !matched {
			return CreateValidatePasswordCriteriaError(ValidatePasswordCriteriaMissingLowerCaseLetter, "password must have at least one lower case letter")
		}
	}

	//validate has upper case letter if required
	if criteria.RequireUpperCase {
		matched, _ := regexp.Match("[A-Z]", []byte(password))
		if !matched {
			return CreateValidatePasswordCriteriaError(ValidatePasswordCriteriaMissingUpperCaseLetter, "password must have at least one upper case letter")
		}
	}

	//validate has digit if required
	if criteria.RequireDigit {
		matched, _ := regexp.Match("\\d", []byte(password))
		if !matched {
			return CreateValidatePasswordCriteriaError(ValidatePasswordCriteriaMissingDigit, "password must have at least one digit")
		}
	}

	//validate has symbol if required
	if criteria.RequireSymbol {
		matched, _ := regexp.Match("[^\\w\\s]", []byte(password))
		if !matched {
			return CreateValidatePasswordCriteriaError(ValidatePasswordCriteriaMissingSymbol, "password must have at least one symbol")
		}
	}

	return CreateValidatePasswordCriteriaValid()
}
