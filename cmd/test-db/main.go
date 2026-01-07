package main

import (
	"flyola-services/internal/config"
	"flyola-services/internal/database"
	"fmt"
	"log"
)

func main() {
	fmt.Println("🧪 Testing Database Connectivity...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Load configuration
	cfg := config.Load()

	fmt.Println("\n📋 Configuration loaded:")
	fmt.Printf("   Database Host: %s\n", cfg.DBHost)
	fmt.Printf("   Database Port: %s\n", cfg.DBPort)
	fmt.Printf("   Database Name: %s\n", cfg.DBName)
	fmt.Printf("   Database User: %s\n", cfg.DBUser)

	// Get DSN
	dsn := cfg.GetDatabaseDSN()
	fmt.Printf("\n🔗 Connection String: %s\n", maskPassword(dsn))

	// Test database connection
	fmt.Println("\n🔌 Attempting to connect to database...")
	db, err := database.Initialize(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	// Test ping
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	// Get database stats
	stats := sqlDB.Stats()
	fmt.Println("\n✅ Database connection successful!")
	fmt.Println("\n📊 Connection Pool Statistics:")
	fmt.Printf("   Open Connections: %d\n", stats.OpenConnections)
	fmt.Printf("   In Use: %d\n", stats.InUse)
	fmt.Printf("   Idle: %d\n", stats.Idle)
	fmt.Printf("   Max Open Connections: %d\n", stats.MaxOpenConnections)

	// Test a simple query
	var version string
	if err := db.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		log.Fatalf("❌ Failed to query database version: %v", err)
	}
	fmt.Printf("\n🗄️  MySQL Version: %s\n", version)

	// List databases
	var databases []string
	if err := db.Raw("SHOW DATABASES").Scan(&databases).Error; err != nil {
		log.Printf("⚠️  Warning: Could not list databases: %v", err)
	} else {
		fmt.Println("\n📚 Available Databases:")
		for _, dbName := range databases {
			marker := " "
			if dbName == cfg.DBName {
				marker = "✓"
			}
			fmt.Printf("   [%s] %s\n", marker, dbName)
		}
	}

	fmt.Println("\n✅ All database connectivity tests passed!")
}

// maskPassword masks the password in the DSN for safe logging
func maskPassword(dsn string) string {
	// Simple masking - find password between : and @
	var masked string
	inPassword := false
	for i, char := range dsn {
		if char == ':' && i > 0 {
			inPassword = true
			masked += string(char)
		} else if char == '@' && inPassword {
			masked += "****" + string(char)
			inPassword = false
		} else if !inPassword {
			masked += string(char)
		}
	}
	return masked
}
