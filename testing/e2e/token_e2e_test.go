package e2e_test

import (
	"net/http"
	"net/url"
	"testing"

	jwthelpers "github.com/mhogar/amber/controllers/jwt_helpers"
	"github.com/mhogar/amber/dependencies"
	"github.com/mhogar/amber/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateTokenRequest(clientID uuid.UUID, username string, password string) *http.Response {
	values := url.Values{
		"client_id": []string{clientID.String()},
		"username":  []string{username},
		"password":  []string{password},
	}
	return suite.SendFormRequest(http.MethodPost, "/token", "", values)
}

type TokenE2ETestSuite struct {
	E2ETestSuite
	User UserCredentials
}

func (suite *TokenE2ETestSuite) SetupSuite() {
	suite.E2ETestSuite.SetupSuite()
	suite.User = suite.CreateUser(suite.AdminToken, "user", 5)
}

func (suite *TokenE2ETestSuite) TearDownSuite() {
	suite.DeleteUser(suite.AdminToken, suite.User.Username)
	suite.E2ETestSuite.TearDownSuite()
}

func (suite *TokenE2ETestSuite) TestCreateToken_UsingDefaultTokenType_RedirectsToURLWithToken() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeDefault, "keys/test.private.pem")

	//create user-role
	role := "role"
	suite.CreateUserRole(suite.AdminToken, clientId, suite.User.Username, "role")

	//create token
	res := suite.SendCreateTokenRequest(clientId, suite.User.Username, suite.User.Password)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	//parse token from url
	claims := suite.parseDefaultTokenClaims("keys/test.public.pem", res.Request.URL.Query().Get("token"))
	suite.Equal(suite.User.Username, claims.Username)
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

func (suite *TokenE2ETestSuite) TestCreateToken_UsingFirebaseTokenType_RedirectsToURLWithToken() {
	keyUri := "keys/firebase-test.json"

	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeFirebase, keyUri)

	//create user-role
	role := "role"
	suite.CreateUserRole(suite.AdminToken, clientId, suite.User.Username, "role")

	//create token
	res := suite.SendCreateTokenRequest(clientId, suite.User.Username, suite.User.Password)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	//parse token from url
	claims := suite.parseFirebaseTokenClaims(keyUri, res.Request.URL.Query().Get("token"))
	suite.Equal(suite.User.Username, claims.UID)
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
