// Package jwt provides a manager for JWT related operations.
package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the structure of JWT claims. It consists of standard registered claims and
// additional UserID which represents the identity of the user.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// JWTManager is a struct that encapsulates the secret used for signing JWT tokens.
type JWTManager struct {
	secret string
}

// NewJWTManager creates a new JWTManager instance with the provided secret.
func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret}
}

// Generate generates a new JWT token with the provided userID.
// The function returns the signed token string or error.
func (m *JWTManager) Generate(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse validates and parses the provided JWT token string.
// It returns the userID from the token claims.
func (m *JWTManager) Parse(tokenString string) (int, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(m.secret), nil
	})

	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
