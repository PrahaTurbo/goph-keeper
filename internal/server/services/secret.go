package services

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/PrahaTurbo/goph-keeper/internal/server/encryption"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
)

// SecretService is an interface that defines methods for handling secret related operations.
type SecretService interface {
	CreateSecret(ctx context.Context, req *models.Secret) error
	GetUserSecrets(ctx context.Context) ([]models.Secret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret) error
	DeleteSecret(ctx context.Context, secretID int) error
}

type secretService struct {
	repo  repository.SecretRepository
	log   *zerolog.Logger
	crypt encryption.Encryption
}

// NewSecretService creates and returns a new SecretService instance.
func NewSecretService(
	repo repository.SecretRepository,
	log *zerolog.Logger,
	crypt encryption.Encryption,
) SecretService {
	return &secretService{
		repo:  repo,
		log:   log,
		crypt: crypt,
	}
}

// CreateSecret creates a new secret for the user.
func (s *secretService) CreateSecret(ctx context.Context, secretModel *models.Secret) error {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to extract user from context")

		return err
	}

	s.crypt.GenerateKey(userID)

	secret := &repository.Secret{
		UserID: userID,
		Type:   secretModel.Type,
	}

	secret.Content, err = s.crypt.Encrypt(secretModel.Content)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to encrypt content")

		return err
	}

	if secretModel.MetaData != "" {
		secret.MetaData, err = s.crypt.Encrypt(secretModel.MetaData)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to encrypt meta data")

			return err
		}
	}

	if err := s.repo.Create(ctx, secret); err != nil {
		s.log.Error().Err(err).Msg("failed to create secret")

		return err
	}

	return nil
}

// GetUserSecrets retrieves all secrets associated with the user.
func (s *secretService) GetUserSecrets(ctx context.Context) ([]models.Secret, error) {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to extract user from context")

		return nil, err
	}

	secrets, err := s.repo.GetUserSecrets(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user secrets")

		return nil, err
	}

	s.crypt.GenerateKey(userID)

	modelSecrets := make([]models.Secret, len(secrets))
	for i := range secrets {
		secret := models.Secret{
			ID:        secrets[i].ID,
			UserID:    secrets[i].UserID,
			Type:      secrets[i].Type,
			CreatedAt: secrets[i].CreatedAt,
		}

		decryptedContent, err := s.crypt.Decrypt(secrets[i].Content)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to decrypt secret content")

			return nil, err
		}

		secret.Content = decryptedContent

		if secrets[i].MetaData != nil {
			decryptedMeta, err := s.crypt.Decrypt(secrets[i].MetaData)
			if err != nil {
				s.log.Error().Err(err).Msg("failed to decrypt secret meta data")

				return nil, err
			}

			secret.MetaData = decryptedMeta
		}

		modelSecrets[i] = secret
	}

	return modelSecrets, nil
}

// UpdateSecret updates the provided secret.
func (s *secretService) UpdateSecret(ctx context.Context, secretModel *models.Secret) error {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to extract user from context")

		return err
	}

	secret := &repository.Secret{
		ID:     secretModel.ID,
		UserID: userID,
		Type:   secretModel.Type,
	}

	secret.Content, err = s.crypt.Encrypt(secretModel.Content)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to encrypt content")

		return err
	}

	if secretModel.MetaData != "" {
		secret.MetaData, err = s.crypt.Encrypt(secretModel.MetaData)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to encrypt meta data")

			return err
		}
	}

	if err := s.repo.UpdateSecret(ctx, secret); err != nil {
		s.log.Error().Err(err).Msg("failed to update secret")

		return err
	}

	return nil
}

// DeleteSecret removes the secret with provided ID.
func (s *secretService) DeleteSecret(ctx context.Context, secretID int) error {
	userID, err := extractUserIDFromCtx(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to extract user from context")

		return err
	}

	if err := s.repo.DeleteSecret(ctx, secretID, userID); err != nil {
		s.log.Error().Err(err).Msg("failed to delete secret")

		return err
	}

	return nil
}
