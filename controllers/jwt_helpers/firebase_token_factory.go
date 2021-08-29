package jwthelpers

import (
	"authserver/common"
	"authserver/config"
	"authserver/loaders"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type FirebaseServiceJSON struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

type FirebaseClaims struct {
	jwt.StandardClaims
	Algorithm string `json:"alg"`
	UID       string `json:"uid"`
}

type FirebaseTokenFactory struct {
	JSONLoader  loaders.JSONLoader
	TokenSigner TokenSigner
}

func (tf FirebaseTokenFactory) CreateToken(keyUri string, _ uuid.UUID, username string) (string, error) {
	var serviceJSON FirebaseServiceJSON

	//load the service json
	err := tf.JSONLoader.Load(keyUri, &serviceJSON)
	if err != nil {
		return "", common.ChainError("error loading service json", err)
	}

	now := time.Now().Unix()

	//fill out the claims as specified by firebase
	//https://firebase.google.com/docs/auth/admin/create-custom-tokens#create_custom_tokens_using_a_third-party_jwt_library
	claims := FirebaseClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    serviceJSON.ClientEmail,
			Subject:   serviceJSON.ClientEmail,
			Audience:  "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
			IssuedAt:  now,
			ExpiresAt: now + config.GetTokenConfig().Lifetime,
		},
		Algorithm: "RS256",
		UID:       username,
	}

	//create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//sign the token
	signedToken, err := tf.TokenSigner.SignToken(token, []byte(serviceJSON.PrivateKey))
	if err != nil {
		return "", common.ChainError("error signing token", err)
	}

	return signedToken, nil
}
