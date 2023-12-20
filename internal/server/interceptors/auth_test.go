package interceptors

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/jwt"
)

type mockHandler func(ctx context.Context, req interface{}) (interface{}, error)

func (m mockHandler) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	return m(ctx, req)
}

func TestAuthInterceptor_UnaryServerInterceptor(t *testing.T) {
	jwtManager := jwt.NewJWTManager("test-secret")
	token, _ := jwtManager.Generate(1)

	testCases := []struct {
		ctx          context.Context
		name         string
		expectedCode codes.Code
	}{
		{
			name: "successful request",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": fmt.Sprintf("bearer %s", token),
			})),
			expectedCode: codes.OK,
		},
		{
			name:         "missing metadata",
			ctx:          context.Background(),
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "no authorization header",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"auth": fmt.Sprintf("bearer %s", token),
			})),
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "the token is not in the correct format",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": fmt.Sprintf("bearer%s", token),
			})),
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "the token is not a Bearer token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": fmt.Sprintf("scheme %s", token),
			})),
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "invalid token",
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": fmt.Sprintf("bearer %s", token[1:]),
			})),
			expectedCode: codes.Unauthenticated,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			})

			a := NewAuthInterceptor(jwtManager)
			info := &grpc.UnaryServerInfo{
				FullMethod: pb.Secret_Create_FullMethodName,
			}

			_, err := a.UnaryServerInterceptor(tt.ctx, "request", info, handler.Handle)

			if status.Code(err) != codes.OK {
				assert.Equal(t, tt.expectedCode.String(), status.Code(err).String())
			}
		})
	}
}
