package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoService(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		userID int
	}{
		{
			name:   "success: encrypt and decrypt",
			userID: 1,
			input:  "test text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cryptoSrvc := NewCryptoService("secret")
			cryptoSrvc.GenerateKey(tt.userID)

			encryptedData, err := cryptoSrvc.Encrypt(tt.input)
			assert.NoError(t, err)

			decryptedData, err := cryptoSrvc.Decrypt(encryptedData)
			assert.NoError(t, err)

			assert.Equal(t, tt.input, decryptedData)
		})
	}
}
