package usecase

import (
	"backend-go/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpirationDuration = 24 * time.Hour
)

type TokenService struct {
	secretKey string
}

func NewTokenService() (*TokenService, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, domain.ErrInternal
	}

	return &TokenService{
		secretKey: secretKey,
	}, nil
}

func (ts *TokenService) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(TokenExpirationDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(ts.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
