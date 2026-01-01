package models

import "time"

// HotelBooking represents a hotel booking
type HotelBooking struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	BookingReference string    `json:"booking_reference" gorm:"unique;not null"`
	UserID           *uint     `json:"user_id"`
	HotelID          uint      `json:"hotel_id"`
	Hotel            Hotel     `json:"hotel" gorm:"foreignKey:HotelID"`
	RoomID           uint      `json:"room_id"`
	Room             Room      `json:"room" gorm:"foreignKey:RoomID"`
	GuestName        string    `json:"guest_name" gorm:"not null"`
	GuestEmail       string    `json:"guest_email" gorm:"not null"`
	GuestPhone       string    `json:"guest_phone" gorm:"not null"`
	CheckInDate      time.Time `json:"check_in_date"`
	CheckOutDate     time.Time `json:"check_out_date"`
	NumberOfNights   int       `json:"number_of_nights"`
	NumberOfGuests   int       `json:"number_of_guests" gorm:"default:1"`
	RoomPrice        float64   `json:"room_price"`
	ExtraPersons     int       `json:"extra_persons" gorm:"default:0"`
	ExtraPersonPrice float64   `json:"extra_person_price" gorm:"default:0"`
	TotalAmount      float64   `json:"total_amount"`
	TaxAmount        float64   `json:"tax_amount" gorm:"default:0"`
	DiscountAmount   float64   `json:"discount_amount" gorm:"default:0"`
	FinalAmount      float64   `json:"final_amount"`
	BookingStatus    string    `json:"booking_status" gorm:"default:pending"`
	PaymentStatus    string    `json:"payment_status" gorm:"default:pending"`
	SpecialRequests  string    `json:"special_requests"`
	PaymentID        string    `json:"payment_id" gorm:"column:payment_id"`
	PaymentMethod    string    `json:"payment_method" gorm:"column:payment_method"`
	BookingDate      time.Time `json:"booking_date" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// HotelGuest represents guests in a booking
type HotelGuest struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	BookingID   uint         `json:"booking_id"`
	Booking     HotelBooking `json:"booking" gorm:"foreignKey:BookingID"`
	Name        string       `json:"name" gorm:"not null"`
	Age         *int         `json:"age"`
	Gender      string       `json:"gender"`
	IDType      string       `json:"id_type"`
	IDNumber    string       `json:"id_number"`
	IsMainGuest bool         `json:"is_main_guest" gorm:"default:false"`
	CreatedAt   time.Time    `json:"created_at"`
}
