package usecase

import (
	"backend-go/internal/domain"
	"context"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) SignUp(ctx context.Context, email, password string) error {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    cleanEmail,
		Password: string(hashedPassword),
	}

	return uc.repo.Create(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := uc.repo.FindByEmail(ctx, cleanEmail)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// 3. Geramos o Token
	// Dica: Em sistemas maiores, essa lógica de JWT ficaria em um "TokenService" separado.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", domain.ErrInternal
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
