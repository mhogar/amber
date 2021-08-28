package jwthelpers

type TokenFactory interface {
	// CreateToken creates a signed JWT using the key loaded from the key uri.
	// Should also include the username in its claims.
	// Returns the token string any errors.
	CreateToken(keyUri string, username string) (string, error)
}
