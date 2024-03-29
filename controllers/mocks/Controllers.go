// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	common "github.com/mhogar/amber/common"
	controllers "github.com/mhogar/amber/controllers"

	mock "github.com/stretchr/testify/mock"

	models "github.com/mhogar/amber/models"

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

// CreateClient provides a mock function with given fields: CRUD, client
func (_m *Controllers) CreateClient(CRUD controllers.ClientControllerCRUD, client *models.Client) common.CustomError {
	ret := _m.Called(CRUD, client)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.ClientControllerCRUD, *models.Client) common.CustomError); ok {
		r0 = rf(CRUD, client)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
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

// CreateTokenRedirectURL provides a mock function with given fields: CRUD, clientId, username, password
func (_m *Controllers) CreateTokenRedirectURL(CRUD controllers.TokenControllerCRUD, clientId uuid.UUID, username string, password string) (string, common.CustomError) {
	ret := _m.Called(CRUD, clientId, username, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(controllers.TokenControllerCRUD, uuid.UUID, string, string) string); ok {
		r0 = rf(CRUD, clientId, username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.TokenControllerCRUD, uuid.UUID, string, string) common.CustomError); ok {
		r1 = rf(CRUD, clientId, username, password)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: CRUD, username, password, rank
func (_m *Controllers) CreateUser(CRUD controllers.UserControllerCRUD, username string, password string, rank int) (*models.User, common.CustomError) {
	ret := _m.Called(CRUD, username, password, rank)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, string, int) *models.User); ok {
		r0 = rf(CRUD, username, password, rank)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserControllerCRUD, string, string, int) common.CustomError); ok {
		r1 = rf(CRUD, username, password, rank)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// CreateUserRole provides a mock function with given fields: CRUD, role
func (_m *Controllers) CreateUserRole(CRUD controllers.UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	ret := _m.Called(CRUD, role)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserRoleControllerCRUD, *models.UserRole) common.CustomError); ok {
		r0 = rf(CRUD, role)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
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

// DeleteAllUserSessions provides a mock function with given fields: CRUD, username
func (_m *Controllers) DeleteAllUserSessions(CRUD controllers.SessionControllerCRUD, username string) common.CustomError {
	ret := _m.Called(CRUD, username)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.SessionControllerCRUD, string) common.CustomError); ok {
		r0 = rf(CRUD, username)
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

// DeleteUserRole provides a mock function with given fields: CRUD, clientUID, username
func (_m *Controllers) DeleteUserRole(CRUD controllers.UserRoleControllerCRUD, clientUID uuid.UUID, username string) common.CustomError {
	ret := _m.Called(CRUD, clientUID, username)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserRoleControllerCRUD, uuid.UUID, string) common.CustomError); ok {
		r0 = rf(CRUD, clientUID, username)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// GetClients provides a mock function with given fields: CRUD
func (_m *Controllers) GetClients(CRUD controllers.ClientControllerCRUD) ([]*models.Client, common.CustomError) {
	ret := _m.Called(CRUD)

	var r0 []*models.Client
	if rf, ok := ret.Get(0).(func(controllers.ClientControllerCRUD) []*models.Client); ok {
		r0 = rf(CRUD)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Client)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.ClientControllerCRUD) common.CustomError); ok {
		r1 = rf(CRUD)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// GetUserRolesWithLesserRankByClientUID provides a mock function with given fields: CRUD, clientUID, rank
func (_m *Controllers) GetUserRolesWithLesserRankByClientUID(CRUD controllers.UserRoleControllerCRUD, clientUID uuid.UUID, rank int) ([]*models.UserRole, common.CustomError) {
	ret := _m.Called(CRUD, clientUID, rank)

	var r0 []*models.UserRole
	if rf, ok := ret.Get(0).(func(controllers.UserRoleControllerCRUD, uuid.UUID, int) []*models.UserRole); ok {
		r0 = rf(CRUD, clientUID, rank)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserRole)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserRoleControllerCRUD, uuid.UUID, int) common.CustomError); ok {
		r1 = rf(CRUD, clientUID, rank)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// GetUsersWithLesserRank provides a mock function with given fields: CRUD, rank
func (_m *Controllers) GetUsersWithLesserRank(CRUD controllers.UserControllerCRUD, rank int) ([]*models.User, common.CustomError) {
	ret := _m.Called(CRUD, rank)

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, int) []*models.User); ok {
		r0 = rf(CRUD, rank)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserControllerCRUD, int) common.CustomError); ok {
		r1 = rf(CRUD, rank)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
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

// UpdateUser provides a mock function with given fields: CRUD, username, rank
func (_m *Controllers) UpdateUser(CRUD controllers.UserControllerCRUD, username string, rank int) (*models.User, common.CustomError) {
	ret := _m.Called(CRUD, username, rank)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, int) *models.User); ok {
		r0 = rf(CRUD, username, rank)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserControllerCRUD, string, int) common.CustomError); ok {
		r1 = rf(CRUD, username, rank)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}

// UpdateUserPassword provides a mock function with given fields: CRUD, username, password
func (_m *Controllers) UpdateUserPassword(CRUD controllers.UserControllerCRUD, username string, password string) common.CustomError {
	ret := _m.Called(CRUD, username, password)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, string) common.CustomError); ok {
		r0 = rf(CRUD, username, password)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// UpdateUserPasswordWithAuth provides a mock function with given fields: CRUD, username, oldPassword, newPassword
func (_m *Controllers) UpdateUserPasswordWithAuth(CRUD controllers.UserControllerCRUD, username string, oldPassword string, newPassword string) common.CustomError {
	ret := _m.Called(CRUD, username, oldPassword, newPassword)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, string, string) common.CustomError); ok {
		r0 = rf(CRUD, username, oldPassword, newPassword)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// UpdateUserRole provides a mock function with given fields: CRUD, role
func (_m *Controllers) UpdateUserRole(CRUD controllers.UserRoleControllerCRUD, role *models.UserRole) common.CustomError {
	ret := _m.Called(CRUD, role)

	var r0 common.CustomError
	if rf, ok := ret.Get(0).(func(controllers.UserRoleControllerCRUD, *models.UserRole) common.CustomError); ok {
		r0 = rf(CRUD, role)
	} else {
		r0 = ret.Get(0).(common.CustomError)
	}

	return r0
}

// VerifyUserRank provides a mock function with given fields: CRUD, username, rank
func (_m *Controllers) VerifyUserRank(CRUD controllers.UserControllerCRUD, username string, rank int) (bool, common.CustomError) {
	ret := _m.Called(CRUD, username, rank)

	var r0 bool
	if rf, ok := ret.Get(0).(func(controllers.UserControllerCRUD, string, int) bool); ok {
		r0 = rf(CRUD, username, rank)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 common.CustomError
	if rf, ok := ret.Get(1).(func(controllers.UserControllerCRUD, string, int) common.CustomError); ok {
		r1 = rf(CRUD, username, rank)
	} else {
		r1 = ret.Get(1).(common.CustomError)
	}

	return r0, r1
}
