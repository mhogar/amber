// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	models "authserver/models"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// Commit provides a mock function with given fields:
func (_m *Transaction) Commit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateClient provides a mock function with given fields: client
func (_m *Transaction) CreateClient(client *models.Client) error {
	ret := _m.Called(client)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Client) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateMigration provides a mock function with given fields: timestamp
func (_m *Transaction) CreateMigration(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: user
func (_m *Transaction) CreateUser(user *models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllOtherUserSessions provides a mock function with given fields: username, tokem
func (_m *Transaction) DeleteAllOtherUserSessions(username string, tokem uuid.UUID) error {
	ret := _m.Called(username, tokem)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) error); ok {
		r0 = rf(username, tokem)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteClient provides a mock function with given fields: uid
func (_m *Transaction) DeleteClient(uid uuid.UUID) (bool, error) {
	ret := _m.Called(uid)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(uid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteMigrationByTimestamp provides a mock function with given fields: timestamp
func (_m *Transaction) DeleteMigrationByTimestamp(timestamp string) error {
	ret := _m.Called(timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSession provides a mock function with given fields: token
func (_m *Transaction) DeleteSession(token uuid.UUID) (bool, error) {
	ret := _m.Called(token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: username
func (_m *Transaction) DeleteUser(username string) (bool, error) {
	ret := _m.Called(username)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClientByUID provides a mock function with given fields: uid
func (_m *Transaction) GetClientByUID(uid uuid.UUID) (*models.Client, error) {
	ret := _m.Called(uid)

	var r0 *models.Client
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Client); ok {
		r0 = rf(uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestTimestamp provides a mock function with given fields:
func (_m *Transaction) GetLatestTimestamp() (string, bool, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMigrationByTimestamp provides a mock function with given fields: timestamp
func (_m *Transaction) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	ret := _m.Called(timestamp)

	var r0 *models.Migration
	if rf, ok := ret.Get(0).(func(string) *models.Migration); ok {
		r0 = rf(timestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Migration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(timestamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSessionByToken provides a mock function with given fields: token
func (_m *Transaction) GetSessionByToken(token uuid.UUID) (*models.Session, error) {
	ret := _m.Called(token)

	var r0 *models.Session
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Session); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: username
func (_m *Transaction) GetUserByUsername(username string) (*models.User, error) {
	ret := _m.Called(username)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserRoleForClient provides a mock function with given fields: clientUID, username
func (_m *Transaction) GetUserRoleForClient(clientUID uuid.UUID, username string) (*models.UserRole, error) {
	ret := _m.Called(clientUID, username)

	var r0 *models.UserRole
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) *models.UserRole); ok {
		r0 = rf(clientUID, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(clientUID, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserRolesForClient provides a mock function with given fields: clientUID
func (_m *Transaction) GetUserRolesForClient(clientUID uuid.UUID) ([]*models.UserRole, error) {
	ret := _m.Called(clientUID)

	var r0 []*models.UserRole
	if rf, ok := ret.Get(0).(func(uuid.UUID) []*models.UserRole); ok {
		r0 = rf(clientUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(clientUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Rollback provides a mock function with given fields:
func (_m *Transaction) Rollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveSession provides a mock function with given fields: session
func (_m *Transaction) SaveSession(session *models.Session) error {
	ret := _m.Called(session)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Session) error); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Setup provides a mock function with given fields:
func (_m *Transaction) Setup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateClient provides a mock function with given fields: client
func (_m *Transaction) UpdateClient(client *models.Client) (bool, error) {
	ret := _m.Called(client)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.Client) bool); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: user
func (_m *Transaction) UpdateUser(user *models.User) (bool, error) {
	ret := _m.Called(user)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.User) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserPassword provides a mock function with given fields: username, hash
func (_m *Transaction) UpdateUserPassword(username string, hash []byte) (bool, error) {
	ret := _m.Called(username, hash)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, []byte) bool); ok {
		r0 = rf(username, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(username, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserRolesForClient provides a mock function with given fields: clientUID, roles
func (_m *Transaction) UpdateUserRolesForClient(clientUID uuid.UUID, roles []*models.UserRole) error {
	ret := _m.Called(clientUID, roles)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, []*models.UserRole) error); ok {
		r0 = rf(clientUID, roles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
