package main

import (
	"backend-go/internal/adapter/handler/http_handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(deps *Dependencies) *gin.Engine {
	r := gin.Default()

	healthHandler := http_handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		api.POST("/register", deps.UserHandler.Register)
		api.POST("/login", deps.UserHandler.Login)
	}

	return r
}
