// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IAuthRepo is an autogenerated mock type for the IAuthRepo type
type IAuthRepo struct {
	mock.Mock
}

// BlacklistToken provides a mock function with given fields: _a0, _a1
func (_m *IAuthRepo) BlacklistToken(_a0 string, _a1 int) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BlacklistToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsTokenBlacklisted provides a mock function with given fields: _a0
func (_m *IAuthRepo) IsTokenBlacklisted(_a0 string) (bool, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for IsTokenBlacklisted")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAuthRepo creates a new instance of IAuthRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthRepo {
	mock := &IAuthRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}