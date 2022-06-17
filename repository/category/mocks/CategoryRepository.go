// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	mock "github.com/stretchr/testify/mock"
)

// CategoryRepository is an autogenerated mock type for the CategoryRepository type
type CategoryRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, ID
func (_m *CategoryRepository) Delete(ctx context.Context, ID string) error {
	ret := _m.Called(ctx, ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx
func (_m *CategoryRepository) FindAll(ctx context.Context) ([]entity.Category, error) {
	ret := _m.Called(ctx)

	var r0 []entity.Category
	if rf, ok := ret.Get(0).(func(context.Context) []entity.Category); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Category)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: ctx, ID
func (_m *CategoryRepository) FindByID(ctx context.Context, ID string) (entity.Category, error) {
	ret := _m.Called(ctx, ID)

	var r0 entity.Category
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.Category); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Get(0).(entity.Category)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, _a1
func (_m *CategoryRepository) Insert(ctx context.Context, _a1 entity.Category) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Category) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, ID, _a2
func (_m *CategoryRepository) Update(ctx context.Context, ID string, _a2 entity.Category) error {
	ret := _m.Called(ctx, ID, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, entity.Category) error); ok {
		r0 = rf(ctx, ID, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewCategoryRepositoryT interface {
	mock.TestingT
	Cleanup(func())
}

// NewCategoryRepository creates a new instance of CategoryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCategoryRepository(t NewCategoryRepositoryT) *CategoryRepository {
	mock := &CategoryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
