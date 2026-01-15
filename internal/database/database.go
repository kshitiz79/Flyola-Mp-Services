package database

import (
	"flyola-services/internal/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	// Configure GORM with custom logger
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Open database connection
	db, err := gorm.Open(mysql.Open(databaseURL), config)
	if err != nil {
		return nil, err
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	// Test database connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("🔌 Database connection pool configured successfully")

	// Auto-migrate holiday package models
	err = db.AutoMigrate(
		&models.HolidayPackage{},
		&models.PackageSchedule{},
		&models.PackageBooking{},
		&models.PackagePassenger{},
		&models.PackageScheduleBooking{},
	)
	if err != nil {
		log.Printf("Warning: Failed to auto-migrate holiday package models: %v", err)
	} else {
		log.Println("✅ Holiday package models migrated successfully")
	}

	// Auto-migrate models (disabled for now to avoid schema conflicts)
	// err = db.AutoMigrate(
	// 	&models.City{},
	// 	&models.Hotel{},
	// 	&models.RoomCategory{},
	// 	&models.Room{},
	// 	&models.RoomAvailability{},
	// 	&models.MealPlan{},
	// 	&models.HotelBooking{},
	// 	&models.HotelGuest{},
	// 	&models.HotelPayment{},
	// 	&models.HotelReview{},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return db, nil
}
