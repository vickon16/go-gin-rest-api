package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vickon16/go-gin-rest-api/internal/env"
)

type CustomClaims struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateLoginToken(userId int64, email string) (string, error) {
	jwtSecret := env.GetEnvString("JWT_SECRET", "some-secret-123456")
	jwtExpirationMinutes := env.GetEnvInt("JWT_EXPIRATION_MINUTES", 10)
	duration := time.Duration(jwtExpirationMinutes) * time.Minute // 10 * 60 // 10 minutes

	// Token expires in minutes
	expirationTime := time.Now().Add(duration)

	claims := CustomClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gin-tutorial",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	jwtSecret := env.GetEnvString("JWT_SECRET", "some-secret-123456")

	tokenFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(jwtSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, tokenFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Optional: check expiration manually
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
