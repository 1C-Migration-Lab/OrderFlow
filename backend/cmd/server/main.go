package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/yourproject/internal/api"
	"github.com/yourproject/internal/repository"
	"github.com/yourproject/internal/service"
)

func main() {
	// Initialize DB connection
	db, err := repository.NewDB("postgres://username:password@localhost:5432/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	repos := repository.NewRepositories(db)

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
