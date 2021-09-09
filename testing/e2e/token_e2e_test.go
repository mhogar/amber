package e2e_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"authserver/dependencies"
	"authserver/models"
	"authserver/router/handlers"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
)

type TokenE2ETestSuite struct {
	E2ETestSuite
	User1 UserCredentials
}

func (suite *TokenE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.User1 = suite.CreateUser(suite.AdminToken, "user1", 5)
}

func (suite *TokenE2ETestSuite) TearDownSuite() {
	suite.DeleteUser(suite.AdminToken, suite.User1.Username)
	suite.E2ETestSuite.TearDownSuite()
}

func (suite *TokenE2ETestSuite) TestCreateDefaultClient_UpdateUserRole_CreateToken_DeleteClient() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeDefault, "keys/test.private.pem")

	//update user role
	role := "role"
	suite.CreateUserRole(suite.User1.Username, clientId, "role")

	//create token
	postTokenBody := handlers.PostTokenBody{
		ClientId: clientId,
		Username: suite.User1.Username,
		Password: suite.User1.Password,
	}
	res := suite.SendRequest(http.MethodPost, "/token", "", postTokenBody)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	claims := suite.parseDefaultTokenClaims("keys/test.public.pem", res.Request.URL.Query().Get("token"))
	suite.Equal(postTokenBody.Username, claims.Username)
	suite.Equal(role, claims.Role)

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) parseDefaultTokenClaims(keyUri string, tokenString string) jwthelpers.DefaultClaims {
	var claims jwthelpers.DefaultClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(_ *jwt.Token) (interface{}, error) {
		//load the public key
		bytes, err := dependencies.ResolveRawDataLoader().Load(keyUri)
		suite.Require().NoError(err)

		//parse the public key
		key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
		suite.Require().NoError(err)

		return key, nil
	})

	suite.Require().NoError(err)
	return claims
}

func (suite *TokenE2ETestSuite) TestCreateFirebaseClient_UpdateUserRole_CreateToken_DeleteClient() {
	keyUri := "keys/firebase-test.json"

	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeFirebase, keyUri)

	//create user role
	role := "role"
	suite.CreateUserRole(suite.User1.Username, clientId, "role")

	//create token
	postTokenBody := handlers.PostTokenBody{
		ClientId: clientId,
		Username: suite.User1.Username,
		Password: suite.User1.Password,
	}
	res := suite.SendRequest(http.MethodPost, "/token", "", postTokenBody)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	claims := suite.parseFirebaseTokenClaims(keyUri, res.Request.URL.Query().Get("token"))
	suite.Equal(suite.User1.Username, claims.UID)
	suite.Equal(role, claims.Claims["role"])

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) parseFirebaseTokenClaims(keyUri string, tokenString string) jwthelpers.FirebaseClaims {
	var claims jwthelpers.FirebaseClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(_ *jwt.Token) (interface{}, error) {
		var serviceJSON jwthelpers.FirebaseServiceJSON

		//load the service json
		err := dependencies.ResolveJSONLoader().Load(keyUri, &serviceJSON)
		suite.Require().NoError(err)

		//parse the private key
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(serviceJSON.PrivateKey))
		suite.Require().NoError(err)

		return &key.PublicKey, nil
	})

	suite.Require().NoError(err)
	return claims
}

func TestTokenE2ETestSuite(t *testing.T) {
	suite.Run(t, &TokenE2ETestSuite{})
}
