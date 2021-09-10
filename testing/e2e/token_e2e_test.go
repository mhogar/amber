package e2e_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"authserver/dependencies"
	"authserver/models"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) SendCreateTokenRequest(clientID uuid.UUID, username string, password string) *http.Response {
	postTokenBody := handlers.PostTokenBody{
		ClientId: clientID,
		Username: username,
		Password: password,
	}
	return suite.SendRequest(http.MethodPost, "/token", "", postTokenBody)
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

func (suite *TokenE2ETestSuite) TestCreateToken_WhereClientNotFound_ReturnsBadRequest() {
	res := suite.SendCreateTokenRequest(uuid.New(), suite.User.Username, suite.User.Password)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "client", "not found")
}

func (suite *TokenE2ETestSuite) TestCreateToken_WhereUserNotFound_ReturnsBadRequest() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, 0, "key.pem")

	//create token
	res := suite.SendCreateTokenRequest(clientId, "DNE", "")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid username and/or password")

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) TestCreateToken_WithIncorrectPassword_ReturnsBadRequest() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, 0, "key.pem")

	//create token
	res := suite.SendCreateTokenRequest(clientId, suite.Admin.Username, "incorrect")
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "invalid username and/or password")

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) TestCreateToken_WhereRoleNotFound_ReturnsBadRequest() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, 0, "key.pem")

	//create token
	res := suite.SendCreateTokenRequest(clientId, suite.User.Username, suite.User.Password)
	helpers.ParseAndAssertErrorResponse(&suite.Suite, res, http.StatusBadRequest, "role for user", "not found")

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) TestCreateToken_WithInvalidKeyURI_ReturnsInternalServerError() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeDefault, "invalid")

	//create user-role
	suite.CreateUserRole(suite.AdminToken, suite.User.Username, clientId, "role")

	//create token
	res := suite.SendCreateTokenRequest(clientId, suite.User.Username, suite.User.Password)
	helpers.ParseAndAssertInternalServerErrorResponse(&suite.Suite, res)

	//delete client
	suite.DeleteClient(suite.AdminToken, clientId)
}

func (suite *TokenE2ETestSuite) TestCreateToken_UsingDefaultTokenType_RedirectsToURLWithToken() {
	//create client
	clientId := suite.CreateClient(suite.AdminToken, models.ClientTokenTypeDefault, "keys/test.private.pem")

	//create user-role
	role := "role"
	suite.CreateUserRole(suite.AdminToken, suite.User.Username, clientId, "role")

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
	suite.CreateUserRole(suite.AdminToken, suite.User.Username, clientId, "role")

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
