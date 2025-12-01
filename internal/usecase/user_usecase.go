package usecase

import (
	"backend-go/internal/domain"
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo         domain.UserRepository
	tokenService domain.TokenService
}

func NewUserUseCase(repo domain.UserRepository, tokenService domain.TokenService) *UserUseCase {
	return &UserUseCase{
		repo:         repo,
		tokenService: tokenService,
	}
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

	tokenString, err := uc.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
