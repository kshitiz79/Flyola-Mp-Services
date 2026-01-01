package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type HotelService struct {
	db *gorm.DB
}

func NewHotelService(db *gorm.DB) *HotelService {
	return &HotelService{db: db}
}

func (s *HotelService) GetAllHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := s.db.Preload("City").Where("status = ?", 0).Find(&hotels).Error
	return hotels, err
}

func (s *HotelService) GetHotelByID(id uint) (*models.Hotel, error) {
	var hotel models.Hotel
	err := s.db.Preload("City").First(&hotel, id).Error
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *HotelService) CreateHotel(hotel *models.Hotel) error {
	return s.db.Create(hotel).Error
}

func (s *HotelService) UpdateHotel(id uint, updates *models.Hotel) (*models.Hotel, error) {
	var hotel models.Hotel
	if err := s.db.First(&hotel, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&hotel).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("City").First(&hotel, id).Error; err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *HotelService) DeleteHotel(id uint) error {
	return s.db.Delete(&models.Hotel{}, id).Error
}

func (s *HotelService) GetHotelsByCity(cityID uint) ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := s.db.Preload("City").Where("city_id = ? AND status = ?", cityID, 0).Find(&hotels).Error
	return hotels, err
}

func (s *HotelService) GetHotelsByStarRating(rating int) ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := s.db.Preload("City").Where("star_rating = ? AND status = ?", rating, 0).Find(&hotels).Error
	return hotels, err
}

func (s *HotelService) SearchHotels(query string) ([]models.Hotel, error) {
	var hotels []models.Hotel
	searchPattern := "%" + query + "%"
	err := s.db.Preload("City").Where("(name LIKE ? OR address LIKE ? OR description LIKE ?) AND status = ?", 
		searchPattern, searchPattern, searchPattern, 0).Find(&hotels).Error
	return hotels, err
}