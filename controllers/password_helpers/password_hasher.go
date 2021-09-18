package passwordhelpers

type PasswordHasher interface {
	// HashPassword hashes the passwords and returns the hash and any errors.
	HashPassword(password string) ([]byte, error)

	// ComparePasswords compares a password hash and a plain text password.
	// Returns nil if equal and any other errors.
	ComparePasswords(hash []byte, password string) error
}
