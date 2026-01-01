package models

import "time"

// HotelPayment represents payment information
type HotelPayment struct {
	ID              uint         `json:"id" gorm:"primaryKey"`
	BookingID       uint         `json:"booking_id"`
	Booking         HotelBooking `json:"booking" gorm:"foreignKey:BookingID"`
	PaymentMethod   string       `json:"payment_method" gorm:"not null"`
	TransactionID   string       `json:"transaction_id"`
	Amount          float64      `json:"amount" gorm:"not null"`
	Currency        string       `json:"currency" gorm:"default:INR"`
	Status          string       `json:"status" gorm:"default:pending"`
	PaymentDate     *time.Time   `json:"payment_date"`
	RefundAmount    float64      `json:"refund_amount" gorm:"default:0"`
	RefundDate      *time.Time   `json:"refund_date"`
	GatewayResponse string       `json:"gateway_response"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}