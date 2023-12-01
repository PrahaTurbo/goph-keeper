package services

import (
	"context"
	"errors"

	"github.com/PrahaTurbo/goph-keeper/internal/server/interceptors"
)

var ErrExtractFromContext = errors.New("cannot extract user id from context")

func extractUserIDFromCtx(ctx context.Context) (int, error) {
	userIDVal := ctx.Value(interceptors.UserIDKey)
	userID, ok := userIDVal.(int)
	if !ok {
		return 0, ErrExtractFromContext
	}

	return userID, nil
}
