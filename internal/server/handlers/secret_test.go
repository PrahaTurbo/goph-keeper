package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/mocks"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

func TestSecretHandler_Create(t *testing.T) {
	log := logger.NewLogger()

	type expected struct {
		response *emptypb.Empty
		err      error
	}

	tests := []struct {
		expected expected
		err      error
		req      *pb.CreateRequest
		name     string
	}{
		{
			name: "success: created secret",
			req: &pb.CreateRequest{
				Type:     1,
				Content:  "test",
				MetaData: "test",
			},
			err: nil,
			expected: expected{
				response: &emptypb.Empty{},
				err:      nil,
			},
		},
		{
			name: "error: failed to create secret",
			req: &pb.CreateRequest{
				Type:     1,
				Content:  "test",
				MetaData: "test",
			},
			err: errors.New("test"),
			expected: expected{
				response: nil,
				err:      status.Errorf(codes.Internal, "failed to create secret"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSecretService := new(mocks.MockSecretService)
			mockSecretService.On("CreateSecret", context.Background(), &models.Secret{
				Type:     tt.req.Type.String(),
				Content:  tt.req.Content,
				MetaData: tt.req.MetaData,
			}).Return(tt.err).Times(1)

			handler := NewSecretHandler(mockSecretService, &log)
			response, err := handler.Create(context.Background(), tt.req)

			assert.Equal(t, tt.expected.response, response)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

func TestSecretHandler_GetSecrets(t *testing.T) {
	log := logger.NewLogger()
	now := time.Now()

	type expected struct {
		response *pb.GetSecretsResponse
		err      error
	}

	tests := []struct {
		expected expected
		prepare  func(s *mocks.MockSecretService)
		name     string
	}{
		{
			name: "success: created secret",
			prepare: func(s *mocks.MockSecretService) {
				s.On("GetUserSecrets", context.Background()).
					Return([]models.Secret{
						{
							ID:        10,
							Type:      pb.SecretType_CREDENTIALS.String(),
							Content:   "test",
							MetaData:  "test",
							CreatedAt: now,
						},
					}, nil).Times(1)
			},
			expected: expected{
				response: &pb.GetSecretsResponse{
					Secrets: []*pb.SecretData{
						{
							Id:        10,
							Type:      pb.SecretType_CREDENTIALS,
							Content:   "test",
							MetaData:  "test",
							CreatedAt: timestamppb.New(now),
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "error: failed to create secret",
			prepare: func(s *mocks.MockSecretService) {
				s.On("GetUserSecrets", context.Background()).
					Return(nil, errors.New("test")).
					Times(1)
			},
			expected: expected{
				response: nil,
				err:      status.Errorf(codes.Internal, "failed to get user secrets"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSecretService := new(mocks.MockSecretService)
			tt.prepare(mockSecretService)

			handler := NewSecretHandler(mockSecretService, &log)
			response, err := handler.GetSecrets(context.Background(), &pb.GetSecretsRequest{})

			assert.Equal(t, tt.expected.response, response)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

func TestSecretHandler_Update(t *testing.T) {
	log := logger.NewLogger()

	type expected struct {
		response *emptypb.Empty
		err      error
	}

	tests := []struct {
		expected expected
		err      error
		req      *pb.UpdateRequest
		name     string
	}{
		{
			name: "success: secret updated",
			req: &pb.UpdateRequest{
				SecretId: 10,
				Type:     1,
				Content:  "test",
				MetaData: "test",
			},
			err: nil,
			expected: expected{
				response: &emptypb.Empty{},
				err:      nil,
			},
		},
		{
			name: "error: failed to create secret",
			req: &pb.UpdateRequest{
				SecretId: 10,
				Type:     1,
				Content:  "test",
				MetaData: "test",
			},
			err: errors.New("test"),
			expected: expected{
				response: nil,
				err:      status.Errorf(codes.Internal, "failed to update secret"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSecretService := new(mocks.MockSecretService)
			mockSecretService.On("UpdateSecret", context.Background(), &models.Secret{
				ID:       int(tt.req.SecretId),
				Type:     tt.req.Type.String(),
				Content:  tt.req.Content,
				MetaData: tt.req.MetaData,
			}).Return(tt.err).Times(1)

			handler := NewSecretHandler(mockSecretService, &log)
			response, err := handler.Update(context.Background(), tt.req)

			assert.Equal(t, tt.expected.response, response)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

func TestSecretHandler_Delete(t *testing.T) {
	log := logger.NewLogger()

	type expected struct {
		response *emptypb.Empty
		err      error
	}

	tests := []struct {
		expected expected
		err      error
		req      *pb.DeleteRequest
		name     string
	}{
		{
			name: "success: secret deleted",
			req:  &pb.DeleteRequest{SecretId: 10},
			err:  nil,
			expected: expected{
				response: &emptypb.Empty{},
				err:      nil,
			},
		},
		{
			name: "error: failed to delete secret",
			req:  &pb.DeleteRequest{SecretId: 10},
			err:  errors.New("test"),
			expected: expected{
				response: nil,
				err:      status.Errorf(codes.Internal, "failed to delete secret"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSecretService := new(mocks.MockSecretService)
			mockSecretService.On("DeleteSecret", context.Background(), int(tt.req.SecretId)).
				Return(tt.err).Times(1)

			handler := NewSecretHandler(mockSecretService, &log)
			response, err := handler.Delete(context.Background(), tt.req)

			assert.Equal(t, tt.expected.response, response)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
