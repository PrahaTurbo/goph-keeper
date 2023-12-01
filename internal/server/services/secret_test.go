package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/interceptors"
	"github.com/PrahaTurbo/goph-keeper/internal/server/mocks"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

type badContextKey struct{}

func Test_secretService_CreateSecret(t *testing.T) {
	log := logger.NewLogger()

	tests := []struct {
		expectedErr       error
		modelsSecret      *models.Secret
		prepareRepo       func(s *mocks.MockSecretRepository)
		prepareEncryption func(e *mocks.MockEncryption)
		name              string
	}{
		{
			name: "success: created secret",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("Create", mock.Anything, &repository.Secret{
					UserID:   1,
					Type:     pb.SecretType_BINARY.String(),
					Content:  []byte("encrypted-data"),
					MetaData: []byte("encrypted-data"),
				}).Return(nil).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return([]byte("encrypted-data"), nil).Times(2)
			},
		},
		{
			name: "error: failed to extract user id from context",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo:       func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {},
			expectedErr:       ErrExtractFromContext,
		},
		{
			name: "error: failed to encrypt",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return(nil, errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		{
			name: "error: failed to encrypt meta data",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "content",
				MetaData: "meta",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", "content").
					Return([]byte("encrypted-data"), nil).Times(1)
				e.On("Encrypt", "meta").
					Return(nil, errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		{
			name: "error: failed to create secret",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("Create", mock.Anything, &repository.Secret{
					UserID:   1,
					Type:     pb.SecretType_BINARY.String(),
					Content:  []byte("encrypted-data"),
					MetaData: []byte("encrypted-data"),
				}).Return(errInternal).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return([]byte("encrypted-data"), nil).Times(2)
			},
			expectedErr: errInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockSecretRepository)
			mockEncryption := new(mocks.MockEncryption)

			tt.prepareRepo(mockRepo)
			tt.prepareEncryption(mockEncryption)

			ctx := context.WithValue(context.Background(), interceptors.UserIDKey, 1)

			if tt.expectedErr == ErrExtractFromContext {
				ctx = context.WithValue(context.Background(), badContextKey{}, 1)
			}

			secretService := NewSecretService(mockRepo, &log, mockEncryption)
			err := secretService.CreateSecret(ctx, tt.modelsSecret)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_secretService_GetUserSecrets(t *testing.T) {
	log := logger.NewLogger()
	now := time.Now()

	type expected struct {
		err     error
		secrets []models.Secret
	}

	tests := []struct {
		expected          expected
		prepareRepo       func(s *mocks.MockSecretRepository)
		prepareEncryption func(e *mocks.MockEncryption)
		name              string
	}{
		{
			name: "success: created secret",
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("GetUserSecrets", mock.Anything, 1).
					Return([]repository.Secret{
						{
							ID:        13,
							UserID:    1,
							Type:      pb.SecretType_BINARY.String(),
							Content:   []byte("encrypted-data"),
							MetaData:  []byte("encrypted-data"),
							CreatedAt: now,
						},
					}, nil).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Decrypt", mock.Anything).
					Return("decrypted-data", nil).Times(2)
			},
			expected: expected{
				secrets: []models.Secret{
					{
						ID:        13,
						UserID:    1,
						Type:      pb.SecretType_BINARY.String(),
						Content:   "decrypted-data",
						MetaData:  "decrypted-data",
						CreatedAt: now,
					},
				},
				err: nil,
			},
		},
		{
			name:              "error: failed to extract user id from context",
			prepareRepo:       func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {},
			expected: expected{
				err: ErrExtractFromContext,
			},
		},
		{
			name: "error: failed to get secrets",
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("GetUserSecrets", mock.Anything, 1).
					Return(nil, errInternal).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {},
			expected: expected{
				err: errInternal,
			},
		},
		{
			name: "success: failed to decrypt content",
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("GetUserSecrets", mock.Anything, 1).
					Return([]repository.Secret{
						{
							ID:        13,
							UserID:    1,
							Type:      pb.SecretType_BINARY.String(),
							Content:   []byte("encrypted-content"),
							MetaData:  []byte("encrypted-data"),
							CreatedAt: now,
						},
					}, nil).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Decrypt", []byte("encrypted-content")).
					Return("", errInternal).Times(1)
			},
			expected: expected{
				err: errInternal,
			},
		},
		{
			name: "success: failed to decrypt meta data",
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("GetUserSecrets", mock.Anything, 1).
					Return([]repository.Secret{
						{
							ID:        13,
							UserID:    1,
							Type:      pb.SecretType_BINARY.String(),
							Content:   []byte("encrypted-content"),
							MetaData:  []byte("encrypted-meta"),
							CreatedAt: now,
						},
					}, nil).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Decrypt", []byte("encrypted-content")).
					Return("decrypted-content", nil).Times(1)
				e.On("Decrypt", []byte("encrypted-meta")).
					Return("", errInternal).Times(1)
			},
			expected: expected{
				err: errInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockSecretRepository)
			mockEncryption := new(mocks.MockEncryption)

			tt.prepareRepo(mockRepo)
			tt.prepareEncryption(mockEncryption)

			ctx := context.WithValue(context.Background(), interceptors.UserIDKey, 1)

			if tt.expected.err == ErrExtractFromContext {
				ctx = context.WithValue(context.Background(), badContextKey{}, 1)
			}

			secretService := NewSecretService(mockRepo, &log, mockEncryption)
			actualSecrets, err := secretService.GetUserSecrets(ctx)

			assert.Equal(t, tt.expected.err, err)
			assert.Equal(t, tt.expected.secrets, actualSecrets)
		})
	}
}

func Test_secretService_UpdateSecret(t *testing.T) {
	log := logger.NewLogger()

	tests := []struct {
		expectedErr       error
		modelsSecret      *models.Secret
		prepareRepo       func(s *mocks.MockSecretRepository)
		prepareEncryption func(e *mocks.MockEncryption)
		name              string
	}{
		{
			name: "success: updated secret",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("UpdateSecret", mock.Anything, &repository.Secret{
					UserID:   1,
					Type:     pb.SecretType_BINARY.String(),
					Content:  []byte("encrypted-data"),
					MetaData: []byte("encrypted-data"),
				}).Return(nil).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return([]byte("encrypted-data"), nil).Times(2)
			},
		},
		{
			name: "error: failed to extract user id from context",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo:       func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {},
			expectedErr:       ErrExtractFromContext,
		},
		{
			name: "error: failed to encrypt",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return(nil, errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		{
			name: "error: failed to encrypt meta data",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "content",
				MetaData: "meta",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", "content").
					Return([]byte("encrypted-data"), nil).Times(1)
				e.On("Encrypt", "meta").
					Return(nil, errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		{
			name: "error: failed to update secret",
			modelsSecret: &models.Secret{
				Type:     pb.SecretType_BINARY.String(),
				Content:  "test",
				MetaData: "test",
			},
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("UpdateSecret", mock.Anything, &repository.Secret{
					UserID:   1,
					Type:     pb.SecretType_BINARY.String(),
					Content:  []byte("encrypted-data"),
					MetaData: []byte("encrypted-data"),
				}).Return(errInternal).Times(1)
			},
			prepareEncryption: func(e *mocks.MockEncryption) {
				e.On("GenerateKey", 1).Times(1)
				e.On("Encrypt", mock.Anything).
					Return([]byte("encrypted-data"), nil).Times(2)
			},
			expectedErr: errInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockSecretRepository)
			mockEncryption := new(mocks.MockEncryption)

			tt.prepareRepo(mockRepo)
			tt.prepareEncryption(mockEncryption)

			ctx := context.WithValue(context.Background(), interceptors.UserIDKey, 1)

			if tt.expectedErr == ErrExtractFromContext {
				ctx = context.WithValue(context.Background(), badContextKey{}, 1)
			}

			secretService := NewSecretService(mockRepo, &log, mockEncryption)
			err := secretService.UpdateSecret(ctx, tt.modelsSecret)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_secretService_DeleteSecret(t *testing.T) {
	log := logger.NewLogger()

	tests := []struct {
		expectedErr error
		prepareRepo func(s *mocks.MockSecretRepository)
		name        string
		secretID    int
	}{
		{
			name:     "success: deleted secret",
			secretID: 132,
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("DeleteSecret", mock.Anything, 132, 1).
					Return(nil).Times(1)
			},
		},
		{
			name:     "error: failed to delete secret",
			secretID: 132,
			prepareRepo: func(s *mocks.MockSecretRepository) {
				s.On("DeleteSecret", mock.Anything, 132, 1).
					Return(errInternal).Times(1)
			},
			expectedErr: errInternal,
		},
		{
			name:        "error: failed to get user id from context",
			secretID:    132,
			prepareRepo: func(s *mocks.MockSecretRepository) {},
			expectedErr: ErrExtractFromContext,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockSecretRepository)
			mockEncryption := new(mocks.MockEncryption)

			tt.prepareRepo(mockRepo)

			ctx := context.WithValue(context.Background(), interceptors.UserIDKey, 1)

			if tt.expectedErr == ErrExtractFromContext {
				ctx = context.WithValue(context.Background(), badContextKey{}, 1)
			}

			secretService := NewSecretService(mockRepo, &log, mockEncryption)
			err := secretService.DeleteSecret(ctx, tt.secretID)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
