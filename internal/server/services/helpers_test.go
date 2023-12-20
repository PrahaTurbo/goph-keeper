package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PrahaTurbo/goph-keeper/internal/server/interceptors"
)

func Test_extractUserIDFromCtx(t *testing.T) {
	type badContextKey string
	var badKey badContextKey = "jwt_token"

	tests := []struct {
		ctx     context.Context
		wantErr error
		name    string
		want    int
	}{
		{
			name: "should return valid user ID",
			ctx:  context.WithValue(context.Background(), interceptors.UserIDKey, 123),
			want: 123,
		},
		{
			name:    "should return error if invalid user ID ",
			ctx:     context.WithValue(context.Background(), interceptors.UserIDKey, "abc"),
			wantErr: ErrExtractFromContext,
		},
		{
			name:    "should return error if invalid key",
			ctx:     context.WithValue(context.Background(), badKey, 123),
			wantErr: ErrExtractFromContext,
		},
		{
			name:    "should return error if missing user ID value",
			ctx:     context.Background(),
			wantErr: ErrExtractFromContext,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractUserIDFromCtx(tt.ctx)

			if err != nil {
				assert.Equal(t, tt.wantErr, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
