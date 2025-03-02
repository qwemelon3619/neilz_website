package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthToken struct {
	UserUUID uuid.UUID `json:"userUUID"`
	jwt.RegisteredClaims
}

// GenerateAccessToken - generate JWT Access Token
func GenerateAccessToken(userUUID uuid.UUID) (string, error) {
	claims := AuthToken{}
	claims.UserUUID = userUUID
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_KEY")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// GenerateRefreshToken - generate JWT Access Token
func GenerateRefreshToken(userUUID uuid.UUID) (string, error) {
	claims := AuthToken{}
	claims.UserUUID = userUUID
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// AuthTokenClaims - Claim Define for JWT..
type AuthTokenClaims struct {
	UserUUID uuid.UUID `json:"userUUID"`
	jwt.RegisteredClaims
}

// ExtractClaimsFromRefreshToken - Extract Claim From Refresh Token
func ExtractClaimsFromRefreshToken(tokenString string) (*AuthTokenClaims, error) {
	claims := &AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
