package main

import (
	"backend-go/internal/adapter/database"
	"backend-go/internal/adapter/handler/http"
	"backend-go/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha crítica ao conectar no banco: %v", err)
	}

	repo := database.NewPostgresRepository(db)
	userUseCase := usecase.NewUserUseCase(repo)
	userHandler := http.NewUserHandler(userUseCase)

	r := gin.Default()

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
