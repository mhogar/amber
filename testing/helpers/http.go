package helpers

import (
	"authserver/common"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/suite"
)

// CreateRequest creates an http request object with the given parameters.
func CreateRequest(suite *suite.Suite, method string, url string, bearerToken string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		bodyReader = bytes.NewReader(bodyStr)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	suite.Require().NoError(err)

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req
}

func CreateDummyRequest(suite *suite.Suite, body interface{}) *http.Request {
	return CreateRequest(suite, "", "", "", body)
}

// ParseResponse parses the provided http response, return its status code and body
func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}

func AssertSuccessResponse(suite *suite.Suite, res interface{}) {
	basicRes := res.(common.BasicResponse)
	suite.True(basicRes.Success)
}

// helpers.ParseAndAssertSuccessResponse asserts the response is a success response
func ParseAndAssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes common.BasicResponse
	status := ParseResponse(suite, res, &basicRes)

	suite.Equal(http.StatusOK, status)
	AssertSuccessResponse(suite, basicRes)
}

func AssertErrorResponse(suite *suite.Suite, res interface{}, expectedErrorSubStrings ...string) {
	errRes := res.(common.ErrorResponse)

	suite.False(errRes.Success)
	AssertContainsSubstrings(suite, errRes.Error, expectedErrorSubStrings...)
}

// ParseAndAssertErrorResponse asserts the response is an error reponse with the expected status and error sub strings
func ParseAndAssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes common.ErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	AssertErrorResponse(suite, errRes)
}

func AssertInternalServerErrorResponse(suite *suite.Suite, res interface{}) {
	AssertErrorResponse(suite, res, "internal error")
}

// ParseAndAssertInternalServerErrorResponse asserts the response is an internal server response
func ParseAndAssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	ParseAndAssertErrorResponse(suite, res, http.StatusInternalServerError, "internal error")
}

func AssertOAuthErrorResponse(suite *suite.Suite, res interface{}, expectedError string, expectedDescriptionSubStrings ...string) {
	errRes := res.(common.OAuthErrorResponse)

	suite.Contains(errRes.Error, expectedError)
	AssertContainsSubstrings(suite, errRes.ErrorDescription, expectedDescriptionSubStrings...)
}

// ParseAndAssertOAuthErrorResponse asserts the response is an oauth error reponse with the expected status, error, and description sub strings
func ParseAndAssertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescriptionSubStrings ...string) {
	var errRes common.OAuthErrorResponse
	status := ParseResponse(suite, res, &errRes)

	suite.Equal(expectedStatus, status)
	AssertOAuthErrorResponse(suite, errRes, expectedError, expectedDescriptionSubStrings...)
}

func AssertAccessTokenResponse(suite *suite.Suite, res interface{}, expectedTokenID string) {
	tokenRes := res.(common.AccessTokenResponse)

	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}

// ParseAndAssertAccessTokenResponse asserts the response is an access token response with the expect token
func ParseAndAssertAccessTokenResponse(suite *suite.Suite, res *http.Response, expectedTokenID string) {
	var tokenRes common.AccessTokenResponse
	status := ParseResponse(suite, res, &tokenRes)

	suite.Equal(http.StatusOK, status)
	AssertAccessTokenResponse(suite, res, expectedTokenID)
}

// ParseAndAssertResponseOK asserts the response has an http OK status and returns the parsed result
func ParseAndAssertResponseOK(suite *suite.Suite, res *http.Response, result interface{}) {
	status := ParseResponse(suite, res, result)
	suite.Equal(http.StatusOK, status)
}
