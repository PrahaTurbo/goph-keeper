// Package interceptors provides the gRPC interceptors for the application.
package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/jwt"
)

type UserIDKeyType string

const (
	bearerSchema   = "bearer"
	authentication = "authorization"
)

const UserIDKey UserIDKeyType = "userID"

var unprotectedPaths = map[string]bool{
	pb.Auth_Login_FullMethodName:    true,
	pb.Auth_Register_FullMethodName: true,
}

// AuthInterceptor structure holds the JWT Manager which will be used to parse the token
// from the context metadata for authenticated services.
type AuthInterceptor struct {
	JWTManager *jwt.JWTManager
}

// NewAuthInterceptor is a constructor function that initializes AuthInterceptor.
func NewAuthInterceptor(jwtManager *jwt.JWTManager) AuthInterceptor {
	return AuthInterceptor{
		JWTManager: jwtManager,
	}
}

// UnaryServerInterceptor is a gRPC unary server interceptor function.
// It intercepts each request and if it is not an unprotected path, it checks for Bearer token and uses JWTManager
// To parse and validate the token. The parsed UserID from the token is then added to context.
// It then calls the underlying gRPC handler and returns its response and error.
func (a *AuthInterceptor) UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, ok := unprotectedPaths[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md[authentication]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	splits := strings.SplitN(authHeader[0], " ", 2)
	if len(splits) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "the token is not in the correct format")
	}

	if strings.ToLower(splits[0]) != bearerSchema {
		return nil, status.Errorf(codes.Unauthenticated, "the token is not a Bearer token")
	}

	tokenString := splits[1]

	userID, err := a.JWTManager.Parse(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "the token is invalid")
	}

	newCtx := context.WithValue(ctx, UserIDKey, userID)

	return handler(newCtx, req)
}
