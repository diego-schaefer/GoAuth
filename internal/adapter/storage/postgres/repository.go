package postgres

import (
	"backend-go/internal/domain"
	"context"
	"log"
	"strings"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	err := db.AutoMigrate(&UserModel{})
	if err != nil {
		log.Fatalf("Erro fatal ao migrar banco de dados: %v", err)
	}

	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := FromDomain(user)

	result := r.db.WithContext(ctx).Create(dbUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return domain.ErrUserAlreadyExists
		}
		return result.Error
	}

	user.ID = dbUser.ID

	return nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser UserModel

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&dbUser)

	if result.Error != nil {
		return nil, result.Error
	}

	return dbUser.ToDomain(), nil
}
