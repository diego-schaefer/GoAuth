package domain

import "context"

type User struct {
	ID       string
	Email    string
	Password string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
