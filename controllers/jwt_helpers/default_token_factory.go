package jwthelpers

import (
	"authserver/common"
	"authserver/loaders"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type DefaultClaims struct {
	jwt.StandardClaims
	Username string
}

type DefaultTokenFactory struct {
	DataLoader  loaders.RawDataLoader
	TokenSigner TokenSigner
}

func (tf DefaultTokenFactory) CreateToken(keyUri string, clientUID uuid.UUID, username string) (string, error) {
	now := time.Now().Unix()

	//load the private key
	privateKey, err := tf.DataLoader.Load(keyUri)
	if err != nil {
		return "", common.ChainError("error loading private key", err)
	}

	//fill out the claims
	claims := DefaultClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "", //TODO: add config item
			Audience:  clientUID.String(),
			IssuedAt:  now,
			ExpiresAt: now + 60, //expires in one minute (TODO: add to config)
		},
		Username: username,
	}

	//create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//sign the token
	signedToken, err := tf.TokenSigner.SignToken(token, []byte(privateKey))
	if err != nil {
		return "", common.ChainError("error signing token", err)
	}

	return signedToken, nil
}
