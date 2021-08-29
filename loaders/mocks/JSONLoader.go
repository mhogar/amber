// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// JSONLoader is an autogenerated mock type for the JSONLoader type
type JSONLoader struct {
	mock.Mock
}

// Load provides a mock function with given fields: uri, v
func (_m *JSONLoader) Load(uri string, v interface{}) error {
	ret := _m.Called(uri, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(uri, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
