package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret}
}

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
