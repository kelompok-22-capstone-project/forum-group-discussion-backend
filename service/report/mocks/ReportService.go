// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	payload "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	mock "github.com/stretchr/testify/mock"

	response "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

// ReportService is an autogenerated mock type for the ReportService type
type ReportService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, accessorUserID, p
func (_m *ReportService) Create(ctx context.Context, accessorUserID string, p payload.CreateReport) (string, error) {
	ret := _m.Called(ctx, accessorUserID, p)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, payload.CreateReport) string); ok {
		r0 = rf(ctx, accessorUserID, p)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, payload.CreateReport) error); ok {
		r1 = rf(ctx, accessorUserID, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, accessorRole, status, page, limit
func (_m *ReportService) GetAll(ctx context.Context, accessorRole string, status string, page uint, limit uint) (response.Pagination[response.Report], error) {
	ret := _m.Called(ctx, accessorRole, status, page, limit)

	var r0 response.Pagination[response.Report]
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uint, uint) response.Pagination[response.Report]); ok {
		r0 = rf(ctx, accessorRole, status, page, limit)
	} else {
		r0 = ret.Get(0).(response.Pagination[response.Report])
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, uint, uint) error); ok {
		r1 = rf(ctx, accessorRole, status, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReportService interface {
	mock.TestingT
	Cleanup(func())
}

// NewReportService creates a new instance of ReportService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReportService(t mockConstructorTestingTNewReportService) *ReportService {
	mock := &ReportService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}