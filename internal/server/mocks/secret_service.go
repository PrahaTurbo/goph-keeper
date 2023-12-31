// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/PrahaTurbo/goph-keeper/internal/server/models"
)

// MockSecretService is an autogenerated mock type for the SecretService type
type MockSecretService struct {
	mock.Mock
}

// CreateSecret provides a mock function with given fields: ctx, req
func (_m *MockSecretService) CreateSecret(ctx context.Context, req *models.Secret) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateSecret")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Secret) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSecret provides a mock function with given fields: ctx, secretID
func (_m *MockSecretService) DeleteSecret(ctx context.Context, secretID int) error {
	ret := _m.Called(ctx, secretID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSecret")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, secretID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserSecrets provides a mock function with given fields: ctx
func (_m *MockSecretService) GetUserSecrets(ctx context.Context) ([]models.Secret, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetUserSecrets")
	}

	var r0 []models.Secret
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Secret, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Secret); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Secret)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSecret provides a mock function with given fields: ctx, secret
func (_m *MockSecretService) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	ret := _m.Called(ctx, secret)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSecret")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Secret) error); ok {
		r0 = rf(ctx, secret)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockSecretService creates a new instance of MockSecretService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSecretService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSecretService {
	mock := &MockSecretService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
