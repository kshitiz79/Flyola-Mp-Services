package config

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	Environment    string
	RazorpayID     string
	RazorpaySecret string
}

func Load() *Config {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "root:MyNewPassword123!@tcp(localhost:3306)/flyola?charset=utf8mb4&parseTime=True&loc=Local"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		RazorpayID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpaySecret: getEnv("RAZORPAY_KEY_SECRET", ""),
	}
	
	// Debug logging (don't log secrets in production)
	if cfg.Environment == "development" {
		log.Printf("🔑 Razorpay Key ID: %s", cfg.RazorpayID)
		if cfg.RazorpaySecret != "" {
			log.Printf("🔑 Razorpay Secret: %s***", cfg.RazorpaySecret[:10])
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
