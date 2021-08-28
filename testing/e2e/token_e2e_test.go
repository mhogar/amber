package e2e_test

import (
	jwthelpers "authserver/controllers/jwt_helpers"
	"authserver/dependencies"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenE2ETestSuite struct {
	E2ETestSuite
}

func (suite *TokenE2ETestSuite) Test_CreateUser_Login_CreateClient_CreateToken_DeleteClient_DeleteUser() {
	//create user
	postUserBody := handlers.PostUserBody{
		Username: "username",
		Password: "Password123!",
	}
	res := suite.SendRequest(http.MethodPost, "/user", "", postUserBody)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//login
	token := suite.Login(postUserBody.Username, postUserBody.Password)

	//create client
	postClientBody := handlers.PostClientBody{
		Name:        "Name",
		RedirectUrl: "https://mhogar.dev",
	}
	res = suite.SendRequest(http.MethodPost, "/client", token, postClientBody)
	id := helpers.ParseDataResponseOK(&suite.Suite, res)["id"].(string)

	//create token
	postTokenBody := handlers.PostTokenBody{
		ClientId: uuid.MustParse(id),
		Username: postUserBody.Username,
		Password: postUserBody.Password,
	}
	res = suite.SendRequest(http.MethodPost, "/token", "", postTokenBody)
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	_, claims := suite.parseFirebaseToken(res.Request.URL.Query().Get("token"))
	suite.Equal(postUserBody.Username, claims.UID)

	//delete client
	res = suite.SendRequest(http.MethodDelete, "/client/"+id, token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)

	//delete user
	res = suite.SendRequest(http.MethodDelete, "/user", token, nil)
	helpers.ParseAndAssertOKSuccessResponse(&suite.Suite, res)
}

func TestTokenE2ETestSuite(t *testing.T) {
	suite.Run(t, &TokenE2ETestSuite{})
}

func (suite *TokenE2ETestSuite) parseFirebaseToken(tokenString string) (*jwt.Token, *jwthelpers.FirebaseClaims) {
	keyUri := "keys/firebase-test.json"

	//resolve loaders
	jsonLoader := dependencies.ResolveJSONLoader()
	keyLoader := dependencies.ResolveRSAKeyLoader()

	claims := &jwthelpers.FirebaseClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		//load the service json
		var serviceJSON jwthelpers.FirebaseServiceJSON
		err := jsonLoader.Load(keyUri, &serviceJSON)
		suite.Require().NoError(err)

		//load the private key
		key, err := keyLoader.LoadPrivateKeyFromBytes([]byte(serviceJSON.PrivateKey))
		suite.Require().NoError(err)

		return &key.PublicKey, nil
	})
	suite.Require().NoError(err)

	return token, claims
}
