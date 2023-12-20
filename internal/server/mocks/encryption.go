// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockEncryption is an autogenerated mock type for the Encryption type
type MockEncryption struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: cipherText
func (_m *MockEncryption) Decrypt(cipherText []byte) (string, error) {
	ret := _m.Called(cipherText)

	if len(ret) == 0 {
		panic("no return value specified for Decrypt")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (string, error)); ok {
		return rf(cipherText)
	}
	if rf, ok := ret.Get(0).(func([]byte) string); ok {
		r0 = rf(cipherText)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(cipherText)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: plainText
func (_m *MockEncryption) Encrypt(plainText string) ([]byte, error) {
	ret := _m.Called(plainText)

	if len(ret) == 0 {
		panic("no return value specified for Encrypt")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(plainText)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(plainText)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(plainText)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateKey provides a mock function with given fields: userID
func (_m *MockEncryption) GenerateKey(userID int) {
	_m.Called(userID)
}

// NewMockEncryption creates a new instance of MockEncryption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEncryption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEncryption {
	mock := &MockEncryption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
