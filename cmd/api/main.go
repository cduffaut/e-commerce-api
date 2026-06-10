package main

import (
	"fmt"
	"log"

	"github.com/cduffaut/e-commerce-api/pkg/config"
	"github.com/cduffaut/e-commerce-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Connect to database
	db := database.NewPostgresPool(cfg.DatabaseURL)
	defer db.Close()

	// Create Gin router
	router := gin.Default()

	// Check that server is running (health route)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
