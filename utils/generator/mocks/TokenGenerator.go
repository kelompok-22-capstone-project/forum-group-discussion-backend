// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"
)

// TokenGenerator is an autogenerated mock type for the TokenGenerator type
type TokenGenerator struct {
	mock.Mock
}

// ExtractToken provides a mock function with given fields: c
func (_m *TokenGenerator) ExtractToken(c echo.Context) (string, string, string, bool) {
	ret := _m.Called(c)

	var r0 string
	if rf, ok := ret.Get(0).(func(echo.Context) string); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(echo.Context) string); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(echo.Context) string); ok {
		r2 = rf(c)
	} else {
		r2 = ret.Get(2).(string)
	}

	var r3 bool
	if rf, ok := ret.Get(3).(func(echo.Context) bool); ok {
		r3 = rf(c)
	} else {
		r3 = ret.Get(3).(bool)
	}

	return r0, r1, r2, r3
}

// GenerateToken provides a mock function with given fields: id, username, role, isActive
func (_m *TokenGenerator) GenerateToken(id string, username string, role string, isActive bool) (string, error) {
	ret := _m.Called(id, username, role, isActive)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, bool) string); ok {
		r0 = rf(id, username, role, isActive)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, bool) error); ok {
		r1 = rf(id, username, role, isActive)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewTokenGeneratorT interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenGenerator creates a new instance of TokenGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenGenerator(t NewTokenGeneratorT) *TokenGenerator {
	mock := &TokenGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
