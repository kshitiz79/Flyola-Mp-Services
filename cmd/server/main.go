package main

import (
	"flyola-services/internal/config"
	"flyola-services/internal/database"
	"flyola-services/internal/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✅ Database connected successfully")

	// Initialize router with dependencies
	r := router.Initialize(db, cfg)

	// Start server
	log.Printf("🏨 Flyola Hotel Services Backend starting on port %s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
