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

// CreateDummyRequest creates an http request object with only the provided body.
func CreateDummyRequest(suite *suite.Suite, body interface{}) *http.Request {
	return CreateRequest(suite, "", "", "", body)
}

// ParseResponse parses the provided http response, asserts its status code, and returns its body.
func ParseResponse(suite *suite.Suite, res *http.Response, expectedStatusCode int, body interface{}) {
	suite.Require().Equal(expectedStatusCode, res.StatusCode)

	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)
}

// AssertSuccessResponse asserts the response is a success response.
func AssertSuccessResponse(suite *suite.Suite, res interface{}) {
	basicRes := res.(common.BasicResponse)
	suite.True(basicRes.Success)
}

// ParseAndAssertSuccessResponse parses the response and asserts it is a success response.
func ParseAndAssertSuccessResponse(suite *suite.Suite, res *http.Response) {
	var basicRes common.BasicResponse
	ParseResponse(suite, res, http.StatusOK, &basicRes)

	AssertSuccessResponse(suite, basicRes)
}

// AssertErrorResponse asserts the response is an error reponse with the expected status and error sub strings.
func AssertErrorResponse(suite *suite.Suite, res interface{}, expectedErrorSubStrings ...string) {
	errRes := res.(common.ErrorResponse)

	suite.False(errRes.Success)
	AssertContainsSubstrings(suite, errRes.Error, expectedErrorSubStrings...)
}

// ParseAndAssertErrorResponse parses the response and asserts it is an error reponse with the expected status and error sub strings.
func ParseAndAssertErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes common.ErrorResponse
	ParseResponse(suite, res, expectedStatus, &errRes)

	AssertErrorResponse(suite, errRes)
}

// AssertInternalServerErrorResponse asserts the response is an internal server response.
func AssertInternalServerErrorResponse(suite *suite.Suite, res interface{}) {
	AssertErrorResponse(suite, res, "internal error")
}

// ParseAndAssertInternalServerErrorResponse parses the response and asserts it is an internal server response.
func ParseAndAssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	ParseAndAssertErrorResponse(suite, res, http.StatusInternalServerError, "internal error")
}

// AssertOAuthErrorResponse asserts the response is an oauth error reponse with the expected status, error, and description sub strings.
func AssertOAuthErrorResponse(suite *suite.Suite, res interface{}, expectedError string, expectedDescriptionSubStrings ...string) {
	errRes := res.(common.OAuthErrorResponse)

	suite.Contains(errRes.Error, expectedError)
	AssertContainsSubstrings(suite, errRes.ErrorDescription, expectedDescriptionSubStrings...)
}

// ParseAndAssertOAuthErrorResponse parses the response and asserts it is an oauth error reponse with the expected status, error, and description sub strings.
func ParseAndAssertOAuthErrorResponse(suite *suite.Suite, res *http.Response, expectedStatus int, expectedError string, expectedDescriptionSubStrings ...string) {
	var errRes common.OAuthErrorResponse
	ParseResponse(suite, res, expectedStatus, &errRes)

	AssertOAuthErrorResponse(suite, errRes, expectedError, expectedDescriptionSubStrings...)
}

// AssertAccessTokenResponse asserts the response is an access token response with the expect token.
func AssertAccessTokenResponse(suite *suite.Suite, res interface{}, expectedTokenID string) {
	tokenRes := res.(common.AccessTokenResponse)

	suite.Equal(expectedTokenID, tokenRes.AccessToken)
	suite.Equal("bearer", tokenRes.TokenType)
}

// AssertSuccessDataResponse asserts the response's data field is equivalent to the expected data.
func AssertSuccessDataResponse(suite *suite.Suite, res interface{}, expectedData interface{}) {
	dataRes := res.(common.DataResponse)

	suite.True(dataRes.Success)
	suite.EqualValues(expectedData, dataRes.Data)
}

// ParseResponseOK asserts the response has an http OK status and returns the parsed result.
func ParseResponseOK(suite *suite.Suite, res *http.Response, result interface{}) {
	ParseResponse(suite, res, http.StatusOK, result)
}

// ParseDataResponseOK asserts the response has an http OK status and returns the parsed result from the data field.
func ParseDataResponseOK(suite *suite.Suite, res *http.Response) map[string]interface{} {
	var dataRes common.DataResponse
	ParseResponse(suite, res, http.StatusOK, &dataRes)

	return (dataRes.Data).(map[string]interface{})
}
