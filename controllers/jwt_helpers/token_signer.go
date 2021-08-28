package jwthelpers

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

type TokenSigner interface {
	//SignToken signs the token using the private key.
	//Returns the signed token string and any errors.
	SignToken(token *jwt.Token, key *rsa.PrivateKey) (string, error)
}

type JWTTokenSigner struct{}

func (JWTTokenSigner) SignToken(token *jwt.Token, key *rsa.PrivateKey) (string, error) {
	return token.SignedString(key)
}
