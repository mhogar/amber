// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	data "github.com/mhogar/amber/data"

	mock "github.com/stretchr/testify/mock"
)

// ScopeFactory is an autogenerated mock type for the ScopeFactory type
type ScopeFactory struct {
	mock.Mock
}

// CreateDataExecutorScope provides a mock function with given fields: _a0
func (_m *ScopeFactory) CreateDataExecutorScope(_a0 func(data.DataExecutor) error) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(data.DataExecutor) error) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTransactionScope provides a mock function with given fields: _a0, _a1
func (_m *ScopeFactory) CreateTransactionScope(_a0 data.DataExecutor, _a1 func(data.Transaction) (bool, error)) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(data.DataExecutor, func(data.Transaction) (bool, error)) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
