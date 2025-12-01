package main

import (
	"backend-go/internal/adapter/handler/http_handler"
	postgresrepo "backend-go/internal/adapter/storage/postgres"
	"backend-go/internal/usecase"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB          *gorm.DB
	UserHandler *http_handler.UserHandler
}

func SetupDependencies(cfg *Config) (*Dependencies, error) {
	db, err := setupDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	repo := postgresrepo.NewPostgresRepository(db)

	tokenService, err := usecase.NewTokenService()
	if err != nil {
		return nil, err
	}

	userUseCase := usecase.NewUserUseCase(repo, tokenService)

	userHandler := http_handler.NewUserHandler(userUseCase)

	return &Dependencies{
		DB:          db,
		UserHandler: userHandler,
	}, nil
}

func setupDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
