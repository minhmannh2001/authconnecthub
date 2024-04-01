// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IRoleUC is an autogenerated mock type for the IRoleUC type
type IRoleUC struct {
	mock.Mock
}

// GetRoleIDByName provides a mock function with given fields: _a0
func (_m *IRoleUC) GetRoleIDByName(_a0 string) (uint, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetRoleIDByName")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (uint, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) uint); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIRoleUC creates a new instance of IRoleUC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIRoleUC(t interface {
	mock.TestingT
	Cleanup(func())
}) *IRoleUC {
	mock := &IRoleUC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}