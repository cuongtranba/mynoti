// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/cuongtranba/mynoti/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// ComicRepository is an autogenerated mock type for the ComicRepository type
type ComicRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *ComicRepository) Delete(_a0 context.Context, _a1 int32) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *ComicRepository) Get(_a0 context.Context, _a1 int32) (*domain.Comic, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Comic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) (*domain.Comic, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32) *domain.Comic); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Comic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0
func (_m *ComicRepository) List(_a0 context.Context) ([]domain.Comic, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []domain.Comic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Comic, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Comic); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Comic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *ComicRepository) Save(_a0 context.Context, _a1 *domain.Comic) (*domain.Comic, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *domain.Comic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Comic) (*domain.Comic, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Comic) *domain.Comic); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Comic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Comic) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewComicRepository creates a new instance of ComicRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewComicRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ComicRepository {
	mock := &ComicRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
