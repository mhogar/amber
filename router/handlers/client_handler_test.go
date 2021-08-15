package handlers_test

import (
	"authserver/common"
	"authserver/models"
	"authserver/router/handlers"
	"authserver/testing/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ClientHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithClientErrorCreatingClient_ReturnsBadRequest() {
	//arrange
	body := handlers.PostClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	message := "create client error"
	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithInternalErrorCreatingClient_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithNoErrors_ReturnsClientData() {
	//arrange
	body := handlers.PostClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	client := models.CreateNewClient("name")
	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(client, common.NoError())

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.ClientDataResponse{
		ID:   client.UID.String(),
		Name: client.Name,
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateClient", &suite.TransactionMock, body.Name)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithErrorParsingId_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, nil)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "client id", "invalid format")
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := helpers.CreateDummyRequest(&suite.Suite, "invalid")
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "invalid json body")
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithClientErrorUpdatingClient_ReturnsBadRequest() {
	//arrange
	body := handlers.PutClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	message := "update client error"
	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithInternalErrorUpdatingClient_ReturnsInternalServerError() {
	//arrange
	body := handlers.PutClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithNoErrors_ReturnsClientData() {
	//arrange
	body := handlers.PostClientBody{
		Name: "name",
	}
	req := helpers.CreateDummyRequest(&suite.Suite, body)

	id := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: id.String(),
		},
	}

	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusOK, status)
	helpers.AssertSuccessDataResponse(&suite.Suite, res, handlers.ClientDataResponse{
		ID:   id.String(),
		Name: body.Name,
	})

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateClient", &suite.TransactionMock, mock.MatchedBy(func(client *models.Client) bool {
		return client.UID == id && client.Name == body.Name
	}))
}

func (suite *ClientHandlerTestSuite) TestDeleteClient_WithErrorParsingId_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, "client id", "invalid format")
}

func (suite *ClientHandlerTestSuite) TestDeleteClient_WithClientErrorDeletingUser_ReturnsBadRequest() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	message := "delete client error"
	suite.ControllersMock.On("DeleteClient", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusBadRequest, status)
	helpers.AssertErrorResponse(&suite.Suite, res, message)
}

func (suite *ClientHandlerTestSuite) TestDeleteClient_WithInternalErrorDeletingUser_ReturnsInternalServerError() {
	//arrange
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("DeleteClient", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusInternalServerError, status)
	helpers.AssertInternalServerErrorResponse(&suite.Suite, res)
}

func (suite *ClientHandlerTestSuite) TestDeleteClient_WithNoErrors_ReturnsSuccess() {
	//arrange
	id := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: id.String(),
		},
	}

	suite.ControllersMock.On("DeleteClient", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.TransactionMock)

	//assert
	suite.Equal(http.StatusOK, status)
	helpers.AssertSuccessResponse(&suite.Suite, res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteClient", &suite.TransactionMock, id)
}

func TestClientHandlerTestSuite(t *testing.T) {
	suite.Run(t, &ClientHandlerTestSuite{})
}
