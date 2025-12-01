package main

import (
	"log"
)

func main() {
	// 1. Load Configuration
	cfg := LoadConfig()
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("❌ Configuração inválida: %v", err)
	}

	// 2. Setup Dependencies
	deps, err := SetupDependencies(cfg)
	if err != nil {
		log.Fatalf("❌ Falha ao configurar dependências: %v", err)
	}

	// 3. Setup Routes
	router := SetupRoutes(deps)

	// 4. Start Server
	log.Printf("🚀 Servidor iniciado na porta %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

