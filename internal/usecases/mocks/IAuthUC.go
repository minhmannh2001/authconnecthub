// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	config "github.com/minhmannh2001/authconnecthub/config"
	dto "github.com/minhmannh2001/authconnecthub/internal/dto"

	entity "github.com/minhmannh2001/authconnecthub/internal/entity"

	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// IAuthUC is an autogenerated mock type for the IAuthUC type
type IAuthUC struct {
	mock.Mock
}

// CheckAndRefreshTokens provides a mock function with given fields: _a0, _a1, _a2
func (_m *IAuthUC) CheckAndRefreshTokens(_a0 string, _a1 string, _a2 *config.Config) (string, string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for CheckAndRefreshTokens")
	}

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, *config.Config) (string, string, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(string, string, *config.Config) string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, *config.Config) string); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string, *config.Config) error); ok {
		r2 = rf(_a0, _a1, _a2)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAuthUC) CreateAccessToken(_a0 entity.User, _a1 int) (string, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateAccessToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.User, int) (string, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(entity.User, int) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(entity.User, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRefreshToken provides a mock function with given fields: _a0, _a1, _a2
func (_m *IAuthUC) CreateRefreshToken(_a0 entity.User, _a1 string, _a2 int) (string, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for CreateRefreshToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.User, string, int) (string, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(entity.User, string, int) string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(entity.User, string, int) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsRefreshTokenValidForAccessToken provides a mock function with given fields: _a0, _a1
func (_m *IAuthUC) IsRefreshTokenValidForAccessToken(_a0 string, _a1 string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for IsRefreshTokenValidForAccessToken")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsTokenBlacklisted provides a mock function with given fields: _a0
func (_m *IAuthUC) IsTokenBlacklisted(_a0 string) (bool, error) {
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

// Login provides a mock function with given fields: _a0, _a1
func (_m *IAuthUC) Login(_a0 *gin.Context, _a1 dto.LoginRequestBody) (*dto.JwtTokens, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 *dto.JwtTokens
	var r1 error
	if rf, ok := ret.Get(0).(func(*gin.Context, dto.LoginRequestBody) (*dto.JwtTokens, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(*gin.Context, dto.LoginRequestBody) *dto.JwtTokens); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.JwtTokens)
		}
	}

	if rf, ok := ret.Get(1).(func(*gin.Context, dto.LoginRequestBody) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: c
func (_m *IAuthUC) Logout(c *gin.Context) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gin.Context) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Register provides a mock function with given fields:
func (_m *IAuthUC) Register() {
	_m.Called()
}

// RetrieveFieldFromJwtToken provides a mock function with given fields: _a0, _a1, _a2
func (_m *IAuthUC) RetrieveFieldFromJwtToken(_a0 string, _a1 string, _a2 bool) (interface{}, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for RetrieveFieldFromJwtToken")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, bool) (interface{}, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(string, string, bool) interface{}); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, bool) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: _a0
func (_m *IAuthUC) ValidateToken(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIAuthUC creates a new instance of IAuthUC. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthUC(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthUC {
	mock := &IAuthUC{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}