// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	payload "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	response "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, p
func (_m *UserService) Login(ctx context.Context, p payload.Login) (response.Login, error) {
	ret := _m.Called(ctx, p)

	var r0 response.Login
	if rf, ok := ret.Get(0).(func(context.Context, payload.Login) response.Login); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Get(0).(response.Login)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, payload.Login) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, p
func (_m *UserService) Register(ctx context.Context, p payload.Register) (string, error) {
	ret := _m.Called(ctx, p)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, payload.Register) string); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, payload.Register) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
