package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string
	GinMode     string

	// Database - individual components
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Database - complete URL (takes precedence if set)
	DatabaseURL string

	// Payment Gateway
	RazorpayID     string
	RazorpaySecret string
}

// GetDatabaseDSN returns the MySQL DSN connection string
func (c *Config) GetDatabaseDSN() string {
	// If DATABASE_URL is set, use it directly
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	// Otherwise, build DSN from individual components
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  Warning: Error loading .env file: %v", err)
		log.Printf("ℹ️  Continuing with system environment variables...")
	} else {
		log.Println("✅ Successfully loaded .env file")
	}

	cfg := &Config{
		// Server Configuration
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),

		// Database Configuration - Individual Components
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "flyola_services"),

		// Database Configuration - Complete URL (optional, takes precedence)
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// Payment Gateway
		RazorpayID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpaySecret: getEnv("RAZORPAY_KEY_SECRET", ""),
	}

	// Debug logging (don't log secrets in production)
	if cfg.Environment == "development" {
		log.Printf("🌍 Environment: %s", cfg.Environment)
		log.Printf("🚀 Server Port: %s", cfg.Port)
		log.Printf("🗄️  Database: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
		log.Printf("🔑 Razorpay Key ID: %s", cfg.RazorpayID)
		if cfg.RazorpaySecret != "" {
			log.Printf("🔑 Razorpay Secret: %s***", cfg.RazorpaySecret[:min(10, len(cfg.RazorpaySecret))])
		} else {
			log.Printf("🔑 Razorpay Secret: NOT SET")
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
