package main

import (
	"backend-go/internal/adapter/database"
	"backend-go/internal/usecase"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDatabaseDSN() string {
	// Tenta pegar DATABASE_URL primeiro
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("ERRO: DATABASE_URL não está configurada. Configure a variável de ambiente DATABASE_URL no Render com a string de conexão PostgreSQL completa do Supabase.")
	}

	// Se DATABASE_URL existe e começa com postgres:// ou postgresql://, usa diretamente
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return dsn
	}

	// Se começa com https://, é provável que seja apenas a URL do Supabase, não a string de conexão
	if strings.HasPrefix(dsn, "https://") || strings.HasPrefix(dsn, "http://") {
		log.Fatal("ERRO: DATABASE_URL está configurada com a URL HTTPS do Supabase, mas precisa ser a string de conexão PostgreSQL completa.\n" +
			"Formato esperado: postgresql://postgres:[SENHA]@[HOST]:5432/postgres?sslmode=require\n" +
			"Você pode encontrar essa string no painel do Supabase em: Settings > Database > Connection string > URI")
	}

	// Tenta construir a partir de variáveis de ambiente separadas como fallback
	log.Printf("AVISO: DATABASE_URL não está no formato postgresql://. Tentando construir a partir de variáveis de ambiente separadas...")

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Se não temos variáveis separadas, dá erro claro
	if host == "" && user == "" && password == "" && dbname == "" {
		log.Fatal("ERRO: DATABASE_URL não está no formato correto e variáveis de ambiente separadas não foram encontradas.\n" +
			"Configure DATABASE_URL com o formato: postgresql://postgres:[SENHA]@[HOST]:5432/postgres?sslmode=require\n" +
			"Ou configure as variáveis: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT")
	}

	// Se temos todas as variáveis necessárias, constrói a DSN
	if user != "" && password != "" && dbname != "" {
		if host == "" {
			host = dsn // Usa DATABASE_URL como host se DB_HOST não estiver definido
		}
		if port == "" {
			port = "5432"
		}
		// Remove https:// se presente no host
		host = strings.TrimPrefix(host, "https://")
		host = strings.TrimPrefix(host, "http://")
		host = strings.TrimSuffix(host, "/")

		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", host, user, password, dbname, port)
	}

	// Se não conseguiu construir, retorna erro
	log.Fatal("ERRO: Não foi possível construir a string de conexão. Configure DATABASE_URL (formato: postgresql://postgres:[SENHA]@[HOST]:5432/postgres?sslmode=require) ou todas as variáveis: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT")
	return ""
}

func main() {
	// 1. Conexão com o Banco (Supabase)
	dsn := getDatabaseDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar no banco: %v", err)
	}

	// 2. Injeção de Dependência Manual
	repo := database.NewPostgresRepository(db)
	useCase := usecase.NewUserUseCase(repo)

	// 3. Configurar Router (Gin é muito usado em Go)
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "JSON inválido. Verifique as aspas e vírgulas."})
			return
		}

		err := useCase.SignUp(body.Email, body.Password)
		if err != nil {
			log.Printf("❌ ERRO CRÍTICO NO BANCO: %v", err)

			c.JSON(500, gin.H{
				"error":   "Erro ao criar usuário",
				"details": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{"message": "Usuário criado com sucesso!"})
	})

	r.POST("/login", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "JSON inválido"})
			return
		}

		token, err := useCase.Login(body.Email, body.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"token":   token,
			"message": "Login realizado com sucesso!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor rodando na porta %s", port)

	// r.Run trava o programa aqui e fica ouvindo requisições
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
