package handlers_test

import (
	"net/http"
	"testing"

	"github.com/mhogar/amber/common"
	"github.com/mhogar/amber/models"
	"github.com/mhogar/amber/router/handlers"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ClientHandlerTestSuite struct {
	HandlersTestSuite
}

func (suite *ClientHandlerTestSuite) TestGetClients_WithClientErrorGettingClients_ReturnsBadRequest() {
	//arrange
	message := "get clients error"
	suite.ControllersMock.On("GetClients", mock.Anything).Return(nil, common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.GetClients(nil, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *ClientHandlerTestSuite) TestGetClients_WithInternalErrorGettingClients_ReturnsInternalServerError() {
	//arrange
	suite.ControllersMock.On("GetClients", mock.Anything).Return(nil, common.InternalError())

	//act
	status, res := suite.CoreHandlers.GetClients(nil, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *ClientHandlerTestSuite) TestGetClients_WithNoErrors_ReturnsClientData() {
	//arrange
	clients := []*models.Client{
		models.CreateNewClient("name1", "redirect1.com", 0, "key1.pem"),
		models.CreateNewClient("name2", "redirect2.com", 1, "key2.pem"),
	}
	suite.ControllersMock.On("GetClients", mock.Anything).Return(clients, common.NoError())

	//act
	status, res := suite.CoreHandlers.GetClients(nil, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, []handlers.ClientDataResponse{
		{
			ID: clients[0].UID.String(),
			PostClientBody: handlers.PostClientBody{
				Name:        clients[0].Name,
				RedirectUrl: clients[0].RedirectUrl,
				TokenType:   clients[0].TokenType,
				KeyUri:      clients[0].KeyUri,
			},
		},
		{
			ID: clients[1].UID.String(),
			PostClientBody: handlers.PostClientBody{
				Name:        clients[1].Name,
				RedirectUrl: clients[1].RedirectUrl,
				TokenType:   clients[1].TokenType,
				KeyUri:      clients[1].KeyUri,
			},
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "GetClients", &suite.CRUDMock)
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := suite.CreateDummyJSONRequest("invalid")

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "invalid json body")
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithClientErrorCreatingClient_ReturnsBadRequest() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	message := "create client error"
	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithInternalErrorCreatingClient_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *ClientHandlerTestSuite) TestPostClient_WithNoErrors_ReturnsClientData() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	var client *models.Client
	suite.ControllersMock.On("CreateClient", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		client = args.Get(1).(*models.Client)
	})

	//act
	status, res := suite.CoreHandlers.PostClient(req, nil, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, handlers.ClientDataResponse{
		ID: client.UID.String(),
		PostClientBody: handlers.PostClientBody{
			Name:        client.Name,
			RedirectUrl: client.RedirectUrl,
			TokenType:   client.TokenType,
			KeyUri:      client.KeyUri,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "CreateClient", &suite.CRUDMock, client)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithErrorParsingId_ReturnsBadRequest() {
	//arrange
	req := suite.CreateDummyJSONRequest(nil)
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: "invalid",
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithInvalidJSONBody_ReturnsBadRequest() {
	//arrange
	req := suite.CreateDummyJSONRequest("invalid")
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "invalid json body")
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithClientErrorUpdatingClient_ReturnsBadRequest() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	message := "update client error"
	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.ClientError(message))

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithInternalErrorUpdatingClient_ReturnsInternalServerError() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.InternalError())

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *ClientHandlerTestSuite) TestPutClient_WithNoErrors_ReturnsClientData() {
	//arrange
	body := handlers.PostClientBody{
		Name:        "name",
		RedirectUrl: "redirect.com",
		TokenType:   0,
		KeyUri:      "key.pem",
	}
	req := suite.CreateDummyJSONRequest(body)

	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uuid.New().String(),
		},
	}

	var client *models.Client
	suite.ControllersMock.On("UpdateClient", mock.Anything, mock.Anything).Return(common.NoError()).Run(func(args mock.Arguments) {
		client = args.Get(1).(*models.Client)
	})

	//act
	status, res := suite.CoreHandlers.PutClient(req, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessDataResponse(res, handlers.ClientDataResponse{
		ID: client.UID.String(),
		PostClientBody: handlers.PostClientBody{
			Name:        client.Name,
			RedirectUrl: client.RedirectUrl,
			TokenType:   client.TokenType,
			KeyUri:      client.KeyUri,
		},
	})

	suite.ControllersMock.AssertCalled(suite.T(), "UpdateClient", &suite.CRUDMock, client)
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
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, "client id", "invalid format")
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
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusBadRequest, status)
	suite.ErrorResponse(res, message)
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
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusInternalServerError, status)
	suite.InternalServerErrorResponse(res)
}

func (suite *ClientHandlerTestSuite) TestDeleteClient_WithNoErrors_ReturnsSuccess() {
	//arrange
	uid := uuid.New()
	params := []httprouter.Param{
		{
			Key:   "id",
			Value: uid.String(),
		},
	}

	suite.ControllersMock.On("DeleteClient", mock.Anything, mock.Anything).Return(common.NoError())

	//act
	status, res := suite.CoreHandlers.DeleteClient(nil, params, nil, &suite.CRUDMock)

	//assert
	suite.Require().Equal(http.StatusOK, status)
	suite.SuccessResponse(res)

	suite.ControllersMock.AssertCalled(suite.T(), "DeleteClient", &suite.CRUDMock, uid)
}

func TestClientHandlerTestSuite(t *testing.T) {
	suite.Run(t, &ClientHandlerTestSuite{})
}
