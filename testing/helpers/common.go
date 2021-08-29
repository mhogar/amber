package helpers

import "github.com/stretchr/testify/suite"

// CreateStringOfLength creates a string of the specified length.
func CreateStringOfLength(length int) string {
	s := ""

	for i := 0; i < length; i++ {
		s += "a"
	}

	return s
}

// AssertContainsSubstrings assets the provided str contains all the expected substrings.
func AssertContainsSubstrings(suite *suite.Suite, str string, expectedSubStrs ...string) {
	for _, expectedSubStr := range expectedSubStrs {
		suite.Contains(str, expectedSubStr)
	}
}
