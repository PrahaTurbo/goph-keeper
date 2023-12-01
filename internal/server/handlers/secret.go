package handlers

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/internal/server/services"
)

// SecretHandler implements the secret-related gRPC service.
type SecretHandler struct {
	pb.UnimplementedSecretServer

	service services.SecretService
	log     *zerolog.Logger
}

// NewSecretHandler is the constructor for the SecretHandler.
func NewSecretHandler(service services.SecretService, log *zerolog.Logger) *SecretHandler {
	return &SecretHandler{
		service: service,
		log:     log,
	}
}

// Create is a gRPC method that allows users to create secrets.
func (h *SecretHandler) Create(ctx context.Context, in *pb.CreateRequest) (*emptypb.Empty, error) {
	secret := models.Secret{
		Type:     in.Type.String(),
		Content:  in.Content,
		MetaData: in.MetaData,
	}

	if err := h.service.CreateSecret(ctx, &secret); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create secret")
	}

	return &emptypb.Empty{}, nil
}

// GetSecrets is a gRPC method that fetches the secrets of a user.
func (h *SecretHandler) GetSecrets(ctx context.Context, in *pb.GetSecretsRequest) (*pb.GetSecretsResponse, error) {
	secrets, err := h.service.GetUserSecrets(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user secrets")
	}

	protoSecrets := make([]*pb.SecretData, len(secrets))
	for i := range secrets {
		secretType := pb.SecretType_UNSPECIFIED
		if v, ok := pb.SecretType_value[secrets[i].Type]; ok {
			secretType = pb.SecretType(v)
		}

		secret := &pb.SecretData{
			Id:        int64(secrets[i].ID),
			Type:      secretType,
			Content:   secrets[i].Content,
			MetaData:  secrets[i].MetaData,
			CreatedAt: timestamppb.New(secrets[i].CreatedAt),
		}

		protoSecrets[i] = secret
	}

	response := pb.GetSecretsResponse{Secrets: protoSecrets}

	return &response, nil
}

// Update is a gRPC method that allows users to update secrets.
func (h *SecretHandler) Update(ctx context.Context, in *pb.UpdateRequest) (*emptypb.Empty, error) {
	secret := &models.Secret{
		ID:       int(in.SecretId),
		Type:     in.Type.String(),
		Content:  in.Content,
		MetaData: in.MetaData,
	}

	if err := h.service.UpdateSecret(ctx, secret); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update secret")
	}

	return &emptypb.Empty{}, nil
}

// Delete is a gRPC method that allows users to delete secrets.
func (h *SecretHandler) Delete(ctx context.Context, in *pb.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.service.DeleteSecret(ctx, int(in.SecretId)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete secret")
	}

	return &emptypb.Empty{}, nil
}
