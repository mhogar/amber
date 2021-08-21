package helpers

import (
	"authserver/common"

	"github.com/stretchr/testify/suite"
)

func AssertNoError(suite *suite.Suite, err common.CustomError) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeNone, err.Type)
}

func AssertClientError(suite *suite.Suite, err common.CustomError, expectedSubStrs ...string) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeClient, err.Type)
	AssertContainsSubstrings(suite, err.Error(), expectedSubStrs...)
}

func AssertInternalError(suite *suite.Suite, err common.CustomError) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeInternal, err.Type)
	AssertContainsSubstrings(suite, err.Error(), "internal error")
}
