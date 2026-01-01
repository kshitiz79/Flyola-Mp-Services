package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
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