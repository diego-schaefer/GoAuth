package domain

import (
	"errors"
	"time"
)

var ErrUserAlreadyExists = errors.New("usuário já cadastrado")

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}
