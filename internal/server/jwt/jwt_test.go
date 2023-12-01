package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var secretKey = "secret"

func TestJWTManager(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		expectError bool
	}{
		{
			name:        "success: correct user ID",
			userID:      123,
			expectError: false,
		},
		{
			name:        "error: invalid token",
			userID:      0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewJWTManager(secretKey)

			token, err := manager.Generate(tt.userID)
			if err != nil {
				t.Fatal(err)
			}

			if tt.userID == 0 {
				token = "invalid.token"
			}

			userID, err := manager.Parse(token)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.userID, userID)
		})
	}
}
