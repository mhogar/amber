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

// AuthenticateUserWithPassword provides a mock function with given fields: CRUD, username, password
func (_m *Controllers) AuthenticateUserWithPassword(CRUD controllers.AuthControllerCRUD, username string, password string) (*models.User, common.CustomError) {
	ret := _m.Called(CRUD, username, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(controllers.AuthControllerCRUD, string, string) *models.User); ok {
		r0 = rf(CRUD, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.AuthControllerCRUD, string, string) common.CustomError); ok {
		r1 = rf(CRUD, username, password)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// CreateClient provides a mock function with given fields: CRUD, name, redirectUrl
func (_m *Controllers) CreateClient(CRUD controllers.ClientControllerCRUD, name string, redirectUrl string) (*models.Client, common.CustomError) {
	ret := _m.Called(CRUD, name, redirectUrl)

	var r0 *models.Client
	if rf, ok := ret.Get(0).(func(controllers.ClientControllerCRUD, string, string) *models.Client); ok {
		r0 = rf(CRUD, name, redirectUrl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Client)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.ClientControllerCRUD, string, string) common.CustomError); ok {
		r1 = rf(CRUD, name, redirectUrl)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// CreateSession provides a mock function with given fields: CRUD, username, password
func (_m *Controllers) CreateSession(CRUD controllers.SessionControllerCRUD, username string, password string) (*models.Session, common.CustomError) {
	ret := _m.Called(CRUD, username, password)

	var r0 *models.Session
	if rf, ok := ret.Get(0).(func(controllers.SessionControllerCRUD, string, string) *models.Session); ok {
		r0 = rf(CRUD, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Session)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.SessionControllerCRUD, string, string) common.CustomError); ok {
		r1 = rf(CRUD, username, password)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// CreateToken provides a mock function with given fields: CRUD, username, password
func (_m *Controllers) CreateToken(CRUD controllers.TokenControllerCRUD, username string, password string) common.CustomError {
	ret := _m.Called(CRUD, username, password)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.TokenControllerCRUD, string, string) common.CustomError); ok {
		r0 = rf(CRUD, username, password)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
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

// DeleteAllOtherUserSessions provides a mock function with given fields: CRUD, username, id
func (_m *Controllers) DeleteAllOtherUserSessions(CRUD controllers.SessionControllerCRUD, username string, id uuid.UUID) common.CustomError {
	ret := _m.Called(CRUD, username, id)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.SessionControllerCRUD, string, uuid.UUID) common.CustomError); ok {
		r0 = rf(CRUD, username, id)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// DeleteClient provides a mock function with given fields: CRUD, uid
func (_m *Controllers) DeleteClient(CRUD controllers.ClientControllerCRUD, uid uuid.UUID) common.CustomError {
	ret := _m.Called(CRUD, uid)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.ClientControllerCRUD, uuid.UUID) common.CustomError); ok {
		r0 = rf(CRUD, uid)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// DeleteSession provides a mock function with given fields: CRUD, id
func (_m *Controllers) DeleteSession(CRUD controllers.SessionControllerCRUD, id uuid.UUID) common.CustomError {
	ret := _m.Called(CRUD, id)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.SessionControllerCRUD, uuid.UUID) common.CustomError); ok {
		r0 = rf(CRUD, id)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: CRUD, username
func (_m *Controllers) DeleteUser(CRUD controllers.UserControllerCRUD, username string) common.CustomError {
	ret := _m.Called(CRUD, username)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string) common.CustomError); ok {
		r0 = rf(CRUD, username)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// UpdateClient provides a mock function with given fields: CRUD, client
func (_m *Controllers) UpdateClient(CRUD controllers.ClientControllerCRUD, client *models.Client) common.CustomError {
	ret := _m.Called(CRUD, client)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.ClientControllerCRUD, *models.Client) common.CustomError); ok {
		r0 = rf(CRUD, client)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// UpdateUserPassword provides a mock function with given fields: CRUD, username, oldPassword, newPassword
func (_m *Controllers) UpdateUserPassword(CRUD controllers.UserControllerCRUD, username string, oldPassword string, newPassword string) common.CustomError {
	ret := _m.Called(CRUD, username, oldPassword, newPassword)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, string, string) common.CustomError); ok {
		r0 = rf(CRUD, username, oldPassword, newPassword)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}
