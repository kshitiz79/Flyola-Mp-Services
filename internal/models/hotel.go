package models

import "time"

// Hotel represents a hotel in the system
type Hotel struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	CityID       uint      `json:"city_id"`
	City         City      `json:"city" gorm:"foreignKey:CityID"`
	Address      string    `json:"address" gorm:"type:text"`
	Description  string    `json:"description" gorm:"type:text"`
	StarRating   int       `json:"star_rating"` // 1-5 stars
	ImageURL     string    `json:"image_url" gorm:"type:varchar(500)"`
	Images       string    `json:"images" gorm:"type:json"`    // JSON array of image objects with Cloudinary URLs
	Amenities    string    `json:"amenities" gorm:"type:json"` // JSON array of amenities
	ContactPhone string    `json:"contact_phone" gorm:"type:varchar(20)"`
	ContactEmail string    `json:"contact_email" gorm:"type:varchar(255)"`
	CheckInTime  string    `json:"check_in_time" gorm:"type:varchar(10)"`  // e.g., "14:00"
	CheckOutTime string    `json:"check_out_time" gorm:"type:varchar(10)"` // e.g., "11:00"
	Status       int       `json:"status"`                                 // 0: Active, 1: Inactive
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}