package jwthelpers

import (
	"github.com/golang-jwt/jwt"
)

type TokenSigner interface {
	// SignToken signs the token using the private key.
	// Returns the signed token string and any errors.
	SignToken(token *jwt.Token, key []byte) (string, error)
}

type JWTTokenSigner struct{}

func (JWTTokenSigner) SignToken(token *jwt.Token, key []byte) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
