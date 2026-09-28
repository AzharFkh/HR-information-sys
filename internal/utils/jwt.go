package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte(os.Getenv("JWT_KEY"))

type JWTClaim struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"roles"`
	jwt.RegisteredClaims
}

type TokenPayload struct {
	UserID uint
	Role   string
}

// digunakan saat login
func GenerateToken(payload TokenPayload) (string, error) {
	expirationTime := time.Now().Add(time.Hour) // 1 jam

	claims := &JWTClaim{
		UserID: payload.UserID,
		Role:   payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JWTKey)
}

// digunakan di midlleware
func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token algorithm")
		}
		return JWTKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
