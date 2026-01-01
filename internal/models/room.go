package models

import (
	"time"

	"gorm.io/datatypes"
)

// RoomCategory represents different types of rooms
type RoomCategory struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	HotelID      *uint          `json:"hotelId" gorm:"column:hotel_id"` // Optional: null for global categories
	Hotel        *Hotel         `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`
	Name         string         `json:"name" gorm:"type:varchar(255);not null"` // Single, Double, Suite, etc.
	Description  string         `json:"description" gorm:"type:text"`
	MaxOccupancy int            `json:"maxOccupancy" gorm:"column:max_occupancy"`          // Maximum number of guests
	BedType      string         `json:"bedType" gorm:"column:bed_type;type:varchar(100)"`  // Single, Double, Queen, King
	RoomSize     string         `json:"roomSize" gorm:"column:room_size;type:varchar(50)"` // e.g., "25 sqm"
	Amenities    datatypes.JSON `json:"amenities" gorm:"type:json"`                        // JSON array of room amenities
	ImageURL     string         `json:"imageUrl" gorm:"column:image_url;type:varchar(500)"`
	Status       int            `json:"status"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at"`
}

// Room represents individual rooms in a hotel
type Room struct {
	ID               uint         `json:"id" gorm:"primaryKey"`
	HotelID          uint         `json:"hotelId" gorm:"column:hotel_id"`
	Hotel            Hotel        `json:"hotel" gorm:"foreignKey:HotelID"`
	RoomCategoryID   uint         `json:"roomCategoryId" gorm:"column:room_category_id"`
	RoomCategory     RoomCategory `json:"roomCategory" gorm:"foreignKey:RoomCategoryID"`
	RoomNumber       string       `json:"roomNumber" gorm:"column:room_number;type:varchar(50)"`
	Floor            int          `json:"floor"`
	BasePrice        float64      `json:"basePrice" gorm:"column:base_price"`                // Base price per night
	SinglePrice      float64      `json:"singlePrice" gorm:"column:single_price"`            // Price for single occupancy
	DoublePrice      float64      `json:"doublePrice" gorm:"column:double_price"`            // Price for double occupancy
	ExtraPersonPrice float64      `json:"extraPersonPrice" gorm:"column:extra_person_price"` // Price per extra person
	MaxExtraPersons  int          `json:"maxExtraPersons" gorm:"column:max_extra_persons"` 
	Status           int          `json:"status"`                                            
	CreatedAt        time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

// RoomAvailability represents room availability for specific dates
type RoomAvailability struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RoomID      uint      `json:"roomId" gorm:"column:room_id"`
	Room        Room      `json:"room" gorm:"foreignKey:RoomID"`
	Date        time.Time `json:"date"`
	IsAvailable bool      `json:"isAvailable" gorm:"column:is_available"`
	Price       float64   `json:"price"` // Dynamic pricing for the date
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName overrides the table name used by RoomAvailability to `room_availability`
func (RoomAvailability) TableName() string {
	return "room_availability"
}
