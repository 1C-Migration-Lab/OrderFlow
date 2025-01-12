package main

import (
	"log"

	"fmt"
	"os"

	"github.com/1C-Migration-Lab/OrderFlow/internal/api"
	"github.com/1C-Migration-Lab/OrderFlow/internal/repository"
	"github.com/1C-Migration-Lab/OrderFlow/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPass := os.Getenv("DB_PASSWORD")
	dbURL := fmt.Sprintf("postgres://postgres.hhdqekkbomdofkinrajy:%s@aws-0-eu-central-1.pooler.supabase.com:6543/postgres", dbPass)
	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := repository.NewPostgresRepository(db)

	// Initialize services
	services := service.NewServices(repos)

	// Initialize router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize API handlers
	api.RegisterRoutes(router, services)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
