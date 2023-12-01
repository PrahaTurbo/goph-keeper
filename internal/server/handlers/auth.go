// Package handlers provides the gRPC implementations for the Auth and Secret service.
package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
	"github.com/PrahaTurbo/goph-keeper/internal/server/services"
)

// AuthHandler implements the auth-related gRPC service.
type AuthHandler struct {
	pb.UnimplementedAuthServer

	service services.AuthService
	log     *zerolog.Logger
}

// NewAuthHandler is the constructor for AuthHandler.
func NewAuthHandler(service services.AuthService, log *zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		log:     log,
	}
}

// Register is a gRPC method that allows users to register to the system.
// It returns the user's token or error.
func (h *AuthHandler) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := h.service.Register(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, "login already exist")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := pb.AuthResponse{Token: token}

	return &resp, nil
}

// Login is a gRPC method that allows users to authenticate themselves.
// It returns the user's token or error.
func (h *AuthHandler) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := h.service.Login(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "login or password is invalid")
	}

	resp := pb.AuthResponse{Token: token}

	return &resp, nil
}
