package services

import (
	"flyola-services/internal/models"
	"time"

	"gorm.io/gorm"
)

type BookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{db: db}
}

func (s *BookingService) GetAllBookings() ([]models.HotelBooking, error) {
	var bookings []models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").Find(&bookings).Error
	return bookings, err
}

func (s *BookingService) GetBookingByID(id uint) (*models.HotelBooking, error) {
	var booking models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *BookingService) CreateBooking(booking *models.HotelBooking) error {
	return s.db.Create(booking).Error
}

func (s *BookingService) UpdateBooking(id uint, updates *models.HotelBooking) (*models.HotelBooking, error) {
	var booking models.HotelBooking
	if err := s.db.First(&booking, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&booking).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Hotel").Preload("Room").First(&booking, id).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) DeleteBooking(id uint) error {
	return s.db.Delete(&models.HotelBooking{}, id).Error
}

func (s *BookingService) CancelBooking(id uint) (*models.HotelBooking, error) {
	var booking models.HotelBooking
	if err := s.db.First(&booking, id).Error; err != nil {
		return nil, err
	}

	booking.BookingStatus = "cancelled"
	if err := s.db.Save(&booking).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) GetBookingsByHotel(hotelID uint) ([]models.HotelBooking, error) {
	var bookings []models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").Where("hotel_id = ?", hotelID).Find(&bookings).Error
	return bookings, err
}

func (s *BookingService) GetBookingsByDateRange(startDate, endDate time.Time) ([]models.HotelBooking, error) {
	var bookings []models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").
		Where("check_in_date BETWEEN ? AND ? OR check_out_date BETWEEN ? AND ?", 
			startDate, endDate, startDate, endDate).Find(&bookings).Error
	return bookings, err
}

func (s *BookingService) GetBookingsByStatus(status string) ([]models.HotelBooking, error) {
	var bookings []models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").Where("booking_status = ?", status).Find(&bookings).Error
	return bookings, err
}

func (s *BookingService) GetBookingsByGuestEmail(email string) ([]models.HotelBooking, error) {
	var bookings []models.HotelBooking
	err := s.db.Preload("Hotel").Preload("Room").Where("guest_email = ?", email).Find(&bookings).Error
	return bookings, err
}