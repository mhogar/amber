// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	rsa "crypto/rsa"

	mock "github.com/stretchr/testify/mock"
)

// RSAKeyLoader is an autogenerated mock type for the RSAKeyLoader type
type RSAKeyLoader struct {
	mock.Mock
}

// LoadPrivateKey provides a mock function with given fields: url
func (_m *RSAKeyLoader) LoadPrivateKey(url string) (*rsa.PrivateKey, error) {
	ret := _m.Called(url)

	var r0 *rsa.PrivateKey
	if rf, ok := ret.Get(0).(func(string) *rsa.PrivateKey); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rsa.PrivateKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
