// Package services provides service-layer logic for authentication and managing secrets.
package services

import (
	"context"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrahaTurbo/goph-keeper/internal/server/jwt"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
)

// AuthService is an interface that defines methods for user registration and login functionalities.
type AuthService interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)
}

type authService struct {
	repo       repository.AuthRepository
	log        *zerolog.Logger
	jwtManager *jwt.JWTManager
}

// NewAuthService creates and returns a new AuthService instance.
func NewAuthService(
	repo repository.AuthRepository,
	log *zerolog.Logger,
	jwtManager *jwt.JWTManager,
) AuthService {
	return &authService{
		repo:       repo,
		log:        log,
		jwtManager: jwtManager,
	}
}

// Register registers a new user with the given login and password, and returns a JWT token.
func (a *authService) Register(ctx context.Context, login string, password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error().Err(err).Str("login", login).Msg("failed to create hash from password")

		return "", err
	}

	user := models.User{
		Login:        login,
		PasswordHash: string(passHash),
	}

	userID, err := a.repo.SaveUser(ctx, user)
	if err != nil {
		a.log.Error().Err(err).Str("login", login).Msg("failed to register user")

		return "", err
	}

	token, err := a.jwtManager.Generate(userID)
	if err != nil {
		return "", err
	}

	a.log.Info().Int("user", userID).Msg("user was created")

	return token, nil
}

// Login checks if the given login and password match a user account, and returns a JWT token.
func (a *authService) Login(ctx context.Context, login string, password string) (string, error) {
	savedUser, err := a.repo.GetUser(ctx, login)
	if err != nil {
		a.log.Error().Err(err).Str("login", login).Msg("cannot find user in database")

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.PasswordHash), []byte(password)); err != nil {
		a.log.Error().Err(err).Str("login", login).Msg("hash and password mismatch")

		return "", err
	}

	token, err := a.jwtManager.Generate(savedUser.ID)
	if err != nil {
		return "", err
	}

	a.log.Info().Int("user", savedUser.ID).Msg("user logged in")

	return token, nil
}
