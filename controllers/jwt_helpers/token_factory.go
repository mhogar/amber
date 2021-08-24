package jwthelpers

import (
	"github.com/golang-jwt/jwt"
)

type customClaims struct {
	jwt.StandardClaims
	Algorithm string `json:"alg"`
	UID       string `json:"uid"`
}

type TokenFactory interface {
	CreateToken(username string) (string, error)
}
