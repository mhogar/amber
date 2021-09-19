package jwthelpers

import (
	"time"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/config"
	"github.com/mhogar/amber/loaders"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type DefaultClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type DefaultTokenFactory struct {
	DataLoader  loaders.RawDataLoader
	TokenSigner TokenSigner
}

func (tf DefaultTokenFactory) CreateToken(keyUri string, clientUID uuid.UUID, username string, role string) (string, error) {
	//load the private key
	privateKey, err := tf.DataLoader.Load(keyUri)
	if err != nil {
		return "", common.ChainError("error loading private key", err)
	}

	now := time.Now().Unix()
	cfg := config.GetTokenConfig()

	//fill out the claims
	claims := DefaultClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cfg.DefaultIssuer,
			Audience:  clientUID.String(),
			IssuedAt:  now,
			ExpiresAt: now + cfg.Lifetime,
		},
		Username: username,
		Role:     role,
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
