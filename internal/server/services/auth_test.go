package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrahaTurbo/goph-keeper/internal/server/jwt"
	"github.com/PrahaTurbo/goph-keeper/internal/server/mocks"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

var errInternal = errors.New("test")

func Test_authService_Register(t *testing.T) {
	log := logger.NewLogger()
	jwtManager := jwt.NewJWTManager("test-secret")

	type expected struct {
		err   error
		token string
	}

	tests := []struct {
		expected expected
		err      error
		name     string
		login    string
		password string
		userID   int
	}{
		{
			name:     "success: user created",
			login:    "test",
			password: "test",
			userID:   1,
			err:      nil,
			expected: expected{
				token: createToken(jwtManager, 1),
				err:   nil,
			},
		},
		{
			name:     "error: failed to create user",
			login:    "test",
			password: "test",
			userID:   0,
			err:      errors.New("test"),
			expected: expected{
				token: "",
				err:   errors.New("test"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			mockRepo.On("SaveUser", context.Background(), mock.Anything).
				Return(tt.userID, tt.err).Times(1)

			authService := NewAuthService(mockRepo, &log, jwtManager)
			token, err := authService.Register(context.Background(), tt.login, tt.password)

			assert.Equal(t, tt.expected.token, token)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

func Test_authService_Login(t *testing.T) {
	log := logger.NewLogger()
	jwtManager := jwt.NewJWTManager("test-secret")

	type expected struct {
		err   error
		token string
	}

	tests := []struct {
		expected expected
		prepare  func(s *mocks.MockAuthRepository)
		name     string
		login    string
		password string
	}{
		{
			name:     "success: user created",
			login:    "login",
			password: "test",
			prepare: func(s *mocks.MockAuthRepository) {
				s.On("GetUser", context.Background(), "login").
					Return(&models.User{
						ID:           1,
						Login:        "login",
						PasswordHash: "$2a$10$lSQ88TSGNM6cR6UAdZWzK.eqUP7GYGk3EmmAzgU5vwFSj5OFnYUKa",
					}, nil).Times(1)
			},
			expected: expected{
				token: createToken(jwtManager, 1),
			},
		},
		{
			name:     "error: password doesn't match",
			login:    "login",
			password: "test2",
			prepare: func(s *mocks.MockAuthRepository) {
				s.On("GetUser", context.Background(), "login").
					Return(&models.User{
						ID:           1,
						Login:        "login",
						PasswordHash: "$2a$10$lSQ88TSGNM6cR6UAdZWzK.eqUP7GYGk3EmmAzgU5vwFSj5OFnYUKa",
					}, nil).Times(1)
			},
			expected: expected{
				err: bcrypt.ErrMismatchedHashAndPassword,
			},
		},
		{
			name:     "error: failed to get user",
			login:    "login",
			password: "password",
			prepare: func(s *mocks.MockAuthRepository) {
				s.On("GetUser", context.Background(), "login").
					Return(nil, errInternal).Times(1)
			},
			expected: expected{
				err: errInternal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAuthRepository)
			tt.prepare(mockRepo)

			authService := NewAuthService(mockRepo, &log, jwtManager)
			token, err := authService.Login(context.Background(), tt.login, tt.password)

			assert.Equal(t, tt.expected.err, err)
			assert.Equal(t, tt.expected.token, token)
		})
	}
}

func createToken(m *jwt.JWTManager, userID int) string {
	token, _ := m.Generate(userID)

	return token
}
