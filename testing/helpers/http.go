package helpers

import (
	"authserver/common"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/stretchr/testify/suite"
)

// CreateRequest creates an http request object with the given parameters and body reader.
func CreateRequest(suite *suite.Suite, method string, url string, bearerToken string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	suite.Require().NoError(err)

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req
}

// CreateJSONRequest creates an http request object with the given parameters and JSON body.
func CreateJSONRequest(suite *suite.Suite, method string, url string, bearerToken string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		bodyReader = bytes.NewReader(bodyStr)
	}

	req := CreateRequest(suite, method, url, bearerToken, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// CreateDummyJSONRequest creates an http request object with only the provided JSON body.
func CreateDummyJSONRequest(suite *suite.Suite, body interface{}) *http.Request {
	return CreateJSONRequest(suite, "", "", "", body)
}

// CreateFormRequest creates an http request object with the given parameters and form body.
func CreateFormRequest(suite *suite.Suite, method string, url string, bearerToken string, body url.Values) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyReader = strings.NewReader(body.Encode())
	}

	req := CreateRequest(suite, method, url, bearerToken, bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// CreateDummyFormRequest creates an http request object with only the provided form body.
func CreateDummyFormRequest(suite *suite.Suite, body url.Values) *http.Request {
	return CreateFormRequest(suite, "POST", "", "", body)
}

// ReadAndAssertRawResponse reads the provided http response, then asserts its status code and its raw body.
func ReadAndAssertRawResponse(suite *suite.Suite, res *http.Response, expectedStatusCode int, expectedData []byte) {
	suite.Require().Equal(expectedStatusCode, res.StatusCode)

	//read the body
	buffer := bytes.Buffer{}
	buffer.ReadFrom(res.Body)

	//assert the data
	suite.Equal(expectedData, buffer.Bytes())
}

// ParseJSONResponse parses the provided http response, asserts its status code, and returns its body.
func ParseJSONResponse(suite *suite.Suite, res *http.Response, expectedStatusCode int, body interface{}) {
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

// ParseAndAssertSuccessResponse parses the response and asserts it has the expected http status and a success body.
func ParseAndAssertSuccessResponse(suite *suite.Suite, expectedStatus int, res *http.Response) {
	var basicRes common.BasicResponse
	ParseJSONResponse(suite, res, expectedStatus, &basicRes)

	AssertSuccessResponse(suite, basicRes)
}

// ParseAndAssertOKSuccessResponse parses the response and asserts it has an OK http status and a success body.
func ParseAndAssertOKSuccessResponse(suite *suite.Suite, res *http.Response) {
	ParseAndAssertSuccessResponse(suite, http.StatusOK, res)
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
	ParseJSONResponse(suite, res, expectedStatus, &errRes)

	AssertErrorResponse(suite, errRes, expectedErrorSubStrings...)
}

// AssertInternalServerErrorResponse asserts the response is an internal server response.
func AssertInternalServerErrorResponse(suite *suite.Suite, res interface{}) {
	AssertErrorResponse(suite, res, "internal error")
}

// ParseAndAssertInternalServerErrorResponse parses the response and asserts it is an internal server response.
func ParseAndAssertInternalServerErrorResponse(suite *suite.Suite, res *http.Response) {
	ParseAndAssertErrorResponse(suite, res, http.StatusInternalServerError, "internal error")
}

// AssertInsufficientPermissionsErrorResponse asserts the response is an insufficient permissions error response.
func AssertInsufficientPermissionsErrorResponse(suite *suite.Suite, res interface{}) {
	AssertErrorResponse(suite, res, "insufficient permissions")
}

// ParseAndAssertInsufficientPermissionsErrorResponse parses the response and asserts it is an insufficient permissions error response.
func ParseAndAssertInsufficientPermissionsErrorResponse(suite *suite.Suite, res *http.Response) {
	ParseAndAssertErrorResponse(suite, res, http.StatusForbidden, "insufficient permissions")
}

// AssertSuccessDataResponse asserts the response's data field is equivalent to the expected data.
func AssertSuccessDataResponse(suite *suite.Suite, res interface{}, expectedData interface{}) {
	dataRes := res.(common.DataResponse)

	suite.True(dataRes.Success)
	suite.EqualValues(expectedData, dataRes.Data)
}

// ParseResponseOK asserts the response has an http OK status and returns the parsed result.
func ParseResponseOK(suite *suite.Suite, res *http.Response, result interface{}) {
	ParseJSONResponse(suite, res, http.StatusOK, result)
}

// ParseDataResponseOK asserts the response has an http OK status and returns the parsed result from the data field.
func ParseDataResponseOK(suite *suite.Suite, res *http.Response) map[string]interface{} {
	var dataRes common.DataResponse
	ParseJSONResponse(suite, res, http.StatusOK, &dataRes)

	return (dataRes.Data).(map[string]interface{})
}
