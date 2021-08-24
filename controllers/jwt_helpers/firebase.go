package jwthelpers

import (
	"authserver/common"
	"authserver/config"
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt"
)

type firebaseServiceJSON struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

type FirebaseTokenFactory struct{}

func (tf FirebaseTokenFactory) CreateToken(username string) (string, error) {
	//parse the service json
	serviceJSON, err := tf.readServiceJSON("mhogar-dev-firebase.json")
	if err != nil {
		return "", common.ChainError("error parsing service json", err)
	}

	now := time.Now().Unix()

	//fill out the claims
	claims := customClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    serviceJSON.ClientEmail,
			Subject:   serviceJSON.ClientEmail,
			Audience:  "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
			IssuedAt:  now,
			ExpiresAt: now + 60, //expires in one min
		},
		Algorithm: "RS256",
		UID:       username,
	}

	//create and sign the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(serviceJSON.PrivateKey)
	if err != nil {
		return "", common.ChainError("error signing token", err)
	}

	return token, nil
}

func (tf FirebaseTokenFactory) readServiceJSON(filename string) (*firebaseServiceJSON, error) {
	var serviceJSON firebaseServiceJSON

	//open the json file
	file, err := os.Open(path.Join(config.GetAppRoot(), "static", "keys", filename))
	if err != nil {
		return nil, common.ChainError("error opening service file", err)
	}
	defer file.Close()

	// decode the json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&serviceJSON)
	if err != nil {
		return nil, common.ChainError("invalid service file json", err)
	}

	return &serviceJSON, nil
}
