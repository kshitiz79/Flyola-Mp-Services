package models

import "time"

// HotelReview represents hotel reviews
type HotelReview struct {
	ID         uint         `json:"id" gorm:"primaryKey"`
	BookingID  uint         `json:"booking_id"`
	Booking    HotelBooking `json:"booking" gorm:"foreignKey:BookingID"`
	HotelID    uint         `json:"hotel_id"`
	Hotel      Hotel        `json:"hotel" gorm:"foreignKey:HotelID"`
	UserID     *uint        `json:"user_id"`
	Rating     int          `json:"rating" gorm:"not null"`
	Title      string       `json:"title"`
	Comment    string       `json:"comment"`
	IsVerified bool         `json:"is_verified" gorm:"default:false"`
	Status     int          `json:"status" gorm:"default:0"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}