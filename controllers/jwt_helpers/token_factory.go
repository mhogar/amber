package jwthelpers

import "github.com/google/uuid"

type TokenFactory interface {
	// CreateToken creates a signed JWT using the key loaded from the key uri.
	// Should also include the username in its claims and optionally the client uid.
	// Returns the token string any errors.
	CreateToken(keyUri string, clientUID uuid.UUID, username string) (string, error)
}
