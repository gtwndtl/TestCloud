package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtWrapper is used to wrap the configuration for JWT
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaims defines the structure of JWT payload
type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token with the given email and returns the signed token string
func (j JwtWrapper) GenerateToken(email string) (string, error) {
	claims := JwtClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.ExpirationHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates the JWT token string and returns the claims if valid
func (j *JwtWrapper) ValidateToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
