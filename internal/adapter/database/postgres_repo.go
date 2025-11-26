package database

import (
	"backend-go/internal/domain"
	"strings"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(user *domain.User) error {
	result := r.db.Create(user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			return domain.ErrUserAlreadyExists
		}
		return result.Error
	}

	return nil
}

func (r *PostgresRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
