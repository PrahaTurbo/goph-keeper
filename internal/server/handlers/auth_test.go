package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/mocks"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

func TestAuthHandler_Register(t *testing.T) {
	log := logger.NewLogger()

	tests := []struct {
		err            error
		expectedErr    error
		expectedOutput *pb.AuthResponse
		name           string
		login          string
		password       string
		token          string
	}{
		{
			name:           "success: user registration",
			login:          "test",
			password:       "12345",
			token:          "token_mock_1",
			err:            nil,
			expectedOutput: &pb.AuthResponse{Token: "token_mock_1"},
			expectedErr:    nil,
		},
		{
			name:           "error: user already exists",
			login:          "test",
			password:       "12345",
			token:          "",
			err:            repository.ErrAlreadyExist,
			expectedOutput: nil,
			expectedErr:    status.Error(codes.AlreadyExists, "login already exist"),
		},
		{
			name:           "error: internal error",
			login:          "test",
			password:       "12345",
			token:          "",
			err:            errors.New("internal error"),
			expectedOutput: nil,
			expectedErr:    status.Error(codes.Internal, "internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(mocks.MockAuthService)
			mockAuthService.On("Register", context.Background(), tt.login, tt.password).
				Return(tt.token, tt.err).
				Times(1)

			handler := NewAuthHandler(mockAuthService, &log)
			output, err := handler.Register(context.Background(), &pb.AuthRequest{
				Login:    tt.login,
				Password: tt.password,
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	log := logger.NewLogger()

	tests := []struct {
		err            error
		expectedErr    error
		expectedOutput *pb.AuthResponse
		name           string
		login          string
		password       string
		token          string
	}{
		{
			name:           "success: user login",
			login:          "test",
			password:       "12345",
			token:          "token_mock_1",
			err:            nil,
			expectedOutput: &pb.AuthResponse{Token: "token_mock_1"},
			expectedErr:    nil,
		},
		{
			name:           "error: invalid password or login",
			login:          "test",
			password:       "12345",
			token:          "",
			err:            errors.New("internal error"),
			expectedOutput: nil,
			expectedErr:    status.Error(codes.Internal, "login or password is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(mocks.MockAuthService)
			mockAuthService.On("Login", context.Background(), tt.login, tt.password).
				Return(tt.token, tt.err).
				Times(1)

			handler := NewAuthHandler(mockAuthService, &log)
			output, err := handler.Login(context.Background(), &pb.AuthRequest{
				Login:    tt.login,
				Password: tt.password,
			})

			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
