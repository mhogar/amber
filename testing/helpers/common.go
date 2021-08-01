package helpers

import "github.com/stretchr/testify/suite"

// AssertContainsSubstrings assets the provided str contains all the expected substrings
func AssertContainsSubstrings(suite *suite.Suite, str string, expectedSubStrs ...string) {
	for _, expectedSubStr := range expectedSubStrs {
		suite.Contains(str, expectedSubStr)
	}
}
