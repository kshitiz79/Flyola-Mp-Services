package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type RoomCategoryService struct {
	db *gorm.DB
}

func NewRoomCategoryService(db *gorm.DB) *RoomCategoryService {
	return &RoomCategoryService{db: db}
}

func (s *RoomCategoryService) GetAllRoomCategories() ([]models.RoomCategory, error) {
	var categories []models.RoomCategory
	err := s.db.Preload("Hotel").Where("status = ?", 0).Find(&categories).Error
	return categories, err
}

func (s *RoomCategoryService) GetRoomCategoryByID(id uint) (*models.RoomCategory, error) {
	var category models.RoomCategory
	err := s.db.Preload("Hotel").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *RoomCategoryService) CreateRoomCategory(category *models.RoomCategory) error {
	return s.db.Create(category).Error
}

func (s *RoomCategoryService) UpdateRoomCategory(id uint, updates *models.RoomCategory) (*models.RoomCategory, error) {
	var category models.RoomCategory
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Hotel").First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *RoomCategoryService) DeleteRoomCategory(id uint) error {
	return s.db.Delete(&models.RoomCategory{}, id).Error
}

func (s *RoomCategoryService) GetRoomCategoriesByHotel(hotelID uint) ([]models.RoomCategory, error) {
	var categories []models.RoomCategory
	err := s.db.Preload("Hotel").Where("hotel_id = ? AND status = ?", hotelID, 0).Find(&categories).Error
	return categories, err
}

func (s *RoomCategoryService) GetGlobalRoomCategories() ([]models.RoomCategory, error) {
	var categories []models.RoomCategory
	err := s.db.Where("hotel_id IS NULL AND status = ?", 0).Find(&categories).Error
	return categories, err
}