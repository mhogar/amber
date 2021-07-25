package router_test

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	datamocks "authserver/data/mocks"
	"authserver/router"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OAuthAuthenticatorTestSuite struct {
	suite.Suite
	CRUDMock           datamocks.DataCRUD
	OAuthAuthenticator router.OAuthAuthenticator
}

func (suite *OAuthAuthenticatorTestSuite) SetupTest() {
	suite.CRUDMock = datamocks.DataCRUD{}
	suite.OAuthAuthenticator = router.OAuthAuthenticator{}
}

func (suite *OAuthAuthenticatorTestSuite) TestAuthenticate_WithNoBearerToken_ReturnsClientRequestError() {
	var req *http.Request

	testCase := func() {
		//act
		token, rerr := suite.OAuthAuthenticator.Authenticate(&suite.CRUDMock, req)

		//assert
		suite.Nil(token)

		common.AssertError(&suite.Suite, rerr, "no bearer token")
		suite.Equal(requesterror.ErrorTypeClient, rerr.Type)
	}

	req = common.CreateRequest(&suite.Suite, "", "", "", nil)
	suite.Run("NoAuthorizationHeader", testCase)

	req.Header.Set("Authorization", "invalid")
	suite.Run("AuthorizationHeaderDoesNotContainBearerToken", testCase)
}

func (suite *OAuthAuthenticatorTestSuite) TestAuthenticate_WithBearerTokenInInvalidFormat_ReturnsClientRequestError() {
	//arrange
	req := common.CreateRequest(&suite.Suite, "", "", "invalid", nil)

	//act
	token, rerr := suite.OAuthAuthenticator.Authenticate(&suite.CRUDMock, req)

	//assert
	suite.Nil(token)

	common.AssertError(&suite.Suite, rerr, "bearer token", "invalid format")
	suite.Equal(requesterror.ErrorTypeClient, rerr.Type)
}

func (suite *OAuthAuthenticatorTestSuite) TestAuthenticate_WithErrorFetchingAccessTokenByID_ReturnsInternalServerRequestError() {
	//arrange
	req := common.CreateRequest(&suite.Suite, "", "", uuid.New().String(), nil)
	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, errors.New(""))

	//act
	token, rerr := suite.OAuthAuthenticator.Authenticate(&suite.CRUDMock, req)

	//assert
	suite.Nil(token)

	common.AssertInternalError(&suite.Suite, rerr)
	suite.Equal(requesterror.ErrorTypeInternal, rerr.Type)
}

func (suite *OAuthAuthenticatorTestSuite) TestAuthenticate_WhereAccessTokenWithIDisNotFound_ReturnsClientRequestError() {
	//arrange
	req := common.CreateRequest(&suite.Suite, "", "", uuid.New().String(), nil)
	suite.CRUDMock.On("GetAccessTokenByID", mock.Anything).Return(nil, nil)

	//act
	token, rerr := suite.OAuthAuthenticator.Authenticate(&suite.CRUDMock, req)

	//assert
	suite.Nil(token)

	common.AssertError(&suite.Suite, rerr, "bearer token", "invalid")
	suite.Equal(requesterror.ErrorTypeClient, rerr.Type)
}

func TestOAuthAuthenticatorTestSuite(t *testing.T) {
	suite.Run(t, &OAuthAuthenticatorTestSuite{})
}
