// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt"

	mock "github.com/stretchr/testify/mock"
)

// TokenSigner is an autogenerated mock type for the TokenSigner type
type TokenSigner struct {
	mock.Mock
}

// SignToken provides a mock function with given fields: token, key
func (_m *TokenSigner) SignToken(token *jwt.Token, key string) (string, error) {
	ret := _m.Called(token, key)

	var r0 string
	if rf, ok := ret.Get(0).(func(*jwt.Token, string) string); ok {
		r0 = rf(token, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jwt.Token, string) error); ok {
		r1 = rf(token, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}