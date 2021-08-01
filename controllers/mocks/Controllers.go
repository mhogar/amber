// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	common "authserver/common"
	controllers "authserver/controllers"

	mock "github.com/stretchr/testify/mock"

	models "authserver/models"

	uuid "github.com/google/uuid"
)

// Controllers is an autogenerated mock type for the Controllers type
type Controllers struct {
	mock.Mock
}

// CreateTokenFromPassword provides a mock function with given fields: CRUD, username, password, clientID, scopeName
func (_m *Controllers) CreateTokenFromPassword(CRUD controllers.TokenControllerCRUD, username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, common.OAuthCustomError) {
	ret := _m.Called(CRUD, username, password, clientID, scopeName)

	var r0 *models.AccessToken
	if rf, ok := ret.Get(0).(func(controllers.TokenControllerCRUD, string, string, uuid.UUID, string) *models.AccessToken); ok {
		r0 = rf(CRUD, username, password, clientID, scopeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AccessToken)
		}
	}

	var r1 common.OAuthCustomError
	if rf, ok := ret.Get(1).(func(controllers.TokenControllerCRUD, string, string, uuid.UUID, string) common.OAuthCustomError); ok {
		r1 = rf(CRUD, username, password, clientID, scopeName)
	} else {
		r1 = ret.Get(1).(common.OAuthCustomError)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: CRUD, username, password
func (_m *Controllers) CreateUser(CRUD controllers.UserControllerCRUD, username string, password string) (*models.User, common.CustomError) {
	ret := _m.Called(CRUD, username, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, string) *models.User); ok {
		r0 = rf(CRUD, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserControllerCRUD, string, string) common.CustomError); ok {
		r1 = rf(CRUD, username, password)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// DeleteAllOtherUserTokens provides a mock function with given fields: CRUD, token
func (_m *Controllers) DeleteAllOtherUserTokens(CRUD controllers.TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	ret := _m.Called(CRUD, token)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.TokenControllerCRUD, *models.AccessToken) common.CustomError); ok {
		r0 = rf(CRUD, token)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// DeleteToken provides a mock function with given fields: CRUD, token
func (_m *Controllers) DeleteToken(CRUD controllers.TokenControllerCRUD, token *models.AccessToken) common.CustomError {
	ret := _m.Called(CRUD, token)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.TokenControllerCRUD, *models.AccessToken) common.CustomError); ok {
		r0 = rf(CRUD, token)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: CRUD, user
func (_m *Controllers) DeleteUser(CRUD controllers.UserControllerCRUD, user *models.User) common.CustomError {
	ret := _m.Called(CRUD, user)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, *models.User) common.CustomError); ok {
		r0 = rf(CRUD, user)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// UpdateUserPassword provides a mock function with given fields: CRUD, user, oldPassword, newPassword
func (_m *Controllers) UpdateUserPassword(CRUD controllers.UserControllerCRUD, user *models.User, oldPassword string, newPassword string) common.CustomError {
	ret := _m.Called(CRUD, user, oldPassword, newPassword)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, *models.User, string, string) common.CustomError); ok {
		r0 = rf(CRUD, user, oldPassword, newPassword)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}