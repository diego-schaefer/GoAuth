package postgres

import (
	"backend-go/internal/domain"
	"time"
)

type UserModel struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (m UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID:       m.ID,
		Email:    m.Email,
		Password: m.Password,
	}
}

func FromDomain(u *domain.User) *UserModel {
	return &UserModel{
		Email:    u.Email,
		Password: u.Password,
	}
}
