package helpers

// CreateStringOfLength creates a string of the specified length.
func CreateStringOfLength(length int) string {
	s := ""

	for i := 0; i < length; i++ {
		s += "a"
	}

	return s
}
