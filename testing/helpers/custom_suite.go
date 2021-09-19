package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mhogar/amber/common"

	"github.com/stretchr/testify/suite"
)

type CustomSuite struct {
	suite.Suite
}

// ContainsSubstrings asserts the provided string contains all the expected substrings.
func (suite *CustomSuite) ContainsSubstrings(str string, expectedSubStrs ...string) {
	for _, expectedSubStr := range expectedSubStrs {
		suite.Contains(str, expectedSubStr)
	}
}

// CustomNoError asserts the provided custom error has type none.
func (suite *CustomSuite) CustomNoError(err common.CustomError) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeNone, err.Type)
}

// CustomClientError asserts the provided custom error has type client and its message contains the all expected sub strings.
func (suite *CustomSuite) CustomClientError(err common.CustomError, expectedSubStrs ...string) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeClient, err.Type)
	suite.ContainsSubstrings(err.Error(), expectedSubStrs...)
}

// CustomInternalError asserts the provided custom error has type internal and an internal error message.
func (suite *CustomSuite) CustomInternalError(err common.CustomError) {
	suite.Require().NotNil(err)
	suite.Equal(common.ErrorTypeInternal, err.Type)
	suite.ContainsSubstrings(err.Error(), "internal error")
}

// CreateRequest creates an http request object with the given parameters and body reader.
func (suite *CustomSuite) CreateRequest(method string, url string, bearerToken string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	suite.Require().NoError(err)

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req
}

// CreateJSONRequest creates an http request object with the given parameters and JSON body.
func (suite *CustomSuite) CreateJSONRequest(method string, url string, bearerToken string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyStr, err := json.Marshal(body)
		suite.Require().NoError(err)

		bodyReader = bytes.NewReader(bodyStr)
	}

	req := suite.CreateRequest(method, url, bearerToken, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// CreateDummyJSONRequest creates an http request object with only the provided JSON body.
func (suite *CustomSuite) CreateDummyJSONRequest(body interface{}) *http.Request {
	return suite.CreateJSONRequest("", "", "", body)
}

// CreateFormRequest creates an http request object with the given parameters and form body.
func (suite *CustomSuite) CreateFormRequest(method string, url string, bearerToken string, body url.Values) *http.Request {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyReader = strings.NewReader(body.Encode())
	}

	req := suite.CreateRequest(method, url, bearerToken, bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// CreateDummyFormRequest creates an http request object with only the provided form body.
func (suite *CustomSuite) CreateDummyFormRequest(body url.Values) *http.Request {
	return suite.CreateFormRequest("POST", "", "", body)
}

// ReadAndAssertRawResponse reads the provided http response, then asserts its status code and its raw body.
func (suite *CustomSuite) ReadAndAssertRawResponse(res *http.Response, expectedStatusCode int, expectedData []byte) {
	suite.Require().Equal(expectedStatusCode, res.StatusCode)

	//read the body
	buffer := bytes.Buffer{}
	buffer.ReadFrom(res.Body)

	//assert the data
	suite.Equal(expectedData, buffer.Bytes())
}

// ParseJSONResponse parses the provided http response, asserts its status code, and returns its body.
func (suite *CustomSuite) ParseJSONResponse(res *http.Response, expectedStatusCode int, body interface{}) {
	suite.Require().Equal(expectedStatusCode, res.StatusCode)

	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)
}

// SuccessResponse asserts the response is a success response.
func (suite *CustomSuite) SuccessResponse(res interface{}) {
	basicRes := res.(common.BasicResponse)
	suite.True(basicRes.Success)
}

// ParseAndAssertSuccessResponse parses the response and asserts it has the expected http status and a success body.
func (suite *CustomSuite) ParseAndAssertSuccessResponse(expectedStatus int, res *http.Response) {
	var basicRes common.BasicResponse
	suite.ParseJSONResponse(res, expectedStatus, &basicRes)

	suite.SuccessResponse(basicRes)
}

// ParseAndAssertOKSuccessResponse parses the response and asserts it has an OK http status and a success body.
func (suite *CustomSuite) ParseAndAssertOKSuccessResponse(res *http.Response) {
	suite.ParseAndAssertSuccessResponse(http.StatusOK, res)
}

// ErrorResponse asserts the response is an error reponse with the expected status and error sub strings.
func (suite *CustomSuite) ErrorResponse(res interface{}, expectedErrorSubStrings ...string) {
	errRes := res.(common.ErrorResponse)

	suite.False(errRes.Success)
	suite.ContainsSubstrings(errRes.Error, expectedErrorSubStrings...)
}

// ParseAndAssertErrorResponse parses the response and asserts it is an error reponse with the expected status and error sub strings.
func (suite *CustomSuite) ParseAndAssertErrorResponse(res *http.Response, expectedStatus int, expectedErrorSubStrings ...string) {
	var errRes common.ErrorResponse
	suite.ParseJSONResponse(res, expectedStatus, &errRes)

	suite.ErrorResponse(errRes, expectedErrorSubStrings...)
}

// InternalServerErrorResponse asserts the response is an internal server response.
func (suite *CustomSuite) InternalServerErrorResponse(res interface{}) {
	suite.ErrorResponse(res, "internal error")
}

// ParseAndAssertInternalServerErrorResponse parses the response and asserts it is an internal server response.
func (suite *CustomSuite) ParseAndAssertInternalServerErrorResponse(res *http.Response) {
	suite.ParseAndAssertErrorResponse(res, http.StatusInternalServerError, "internal error")
}

// InsufficientPermissionsErrorResponse asserts the response is an insufficient permissions error response.
func (suite *CustomSuite) InsufficientPermissionsErrorResponse(res interface{}) {
	suite.ErrorResponse(res, "insufficient permissions")
}

// ParseAndAssertInsufficientPermissionsErrorResponse parses the response and asserts it is an insufficient permissions error response.
func (suite *CustomSuite) ParseAndAssertInsufficientPermissionsErrorResponse(res *http.Response) {
	suite.ParseAndAssertErrorResponse(res, http.StatusForbidden, "insufficient permissions")
}

// SuccessDataResponse asserts the response's data field is equivalent to the expected data.
func (suite *CustomSuite) SuccessDataResponse(res interface{}, expectedData interface{}) {
	dataRes := res.(common.DataResponse)

	suite.True(dataRes.Success)
	suite.EqualValues(expectedData, dataRes.Data)
}

// ParseResponseOK asserts the response has an http OK status and returns the parsed result.
func (suite *CustomSuite) ParseResponseOK(res *http.Response, result interface{}) {
	suite.ParseJSONResponse(res, http.StatusOK, result)
}

// ParseDataResponseOK asserts the response has an http OK status and returns the parsed result from the data field.
func (suite *CustomSuite) ParseDataResponseOK(res *http.Response) map[string]interface{} {
	var dataRes common.DataResponse
	suite.ParseJSONResponse(res, http.StatusOK, &dataRes)

	return (dataRes.Data).(map[string]interface{})
}
