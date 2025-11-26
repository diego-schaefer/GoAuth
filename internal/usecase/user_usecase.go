package usecase

import (
	"backend-go/internal/domain"
	"errors"
	"fmt"
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

func (uc *UserUseCase) SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return uc.repo.Create(user)
}

func (uc *UserUseCase) Login(email, password string) (string, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	fmt.Printf("🔍 TENTATIVA DE LOGIN:\n")
	fmt.Printf("   Email Recebido (Limpo): '%s'\n", cleanEmail)

	user, err := uc.repo.FindByEmail(cleanEmail)
	if err != nil {
		fmt.Printf("❌ FALHA NO BANCO: %v\n", err)
		return "", errors.New("email não encontrado no sistema")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("❌ SENHA ERRADA para o email %s\n", cleanEmail)
		return "", errors.New("senha incorreta")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("servidor não configurado (falta JWT_SECRET)")
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
