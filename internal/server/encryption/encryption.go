package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"strconv"

	"golang.org/x/crypto/pbkdf2"
)

type Encryption interface {
	GenerateKey(userID int)
	Encrypt(plainText string) ([]byte, error)
	Decrypt(cipherText []byte) (string, error)
}

type cryptoService struct {
	secret string
	key    []byte
}

func NewCryptoService(secret string) Encryption {
	return &cryptoService{secret: secret}
}

func (e *cryptoService) GenerateKey(userID int) {
	salt := []byte(strconv.Itoa(userID))
	key := pbkdf2.Key([]byte(e.secret), salt, 4096, 32, sha256.New)
	e.key = key
}

func (e *cryptoService) Encrypt(plainText string) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plainText), nil)

	return ciphertext, nil
}

func (e *cryptoService) Decrypt(cipherText []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("cipherText too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
