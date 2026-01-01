package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type CityService struct {
	db *gorm.DB
}

func NewCityService(db *gorm.DB) *CityService {
	return &CityService{db: db}
}

func (s *CityService) GetAllCities() ([]models.City, error) {
	var cities []models.City
	err := s.db.Where("status = ?", 0).Find(&cities).Error
	return cities, err
}

func (s *CityService) GetCityByID(id uint) (*models.City, error) {
	var city models.City
	err := s.db.First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (s *CityService) CreateCity(city *models.City) error {
	return s.db.Create(city).Error
}

func (s *CityService) UpdateCity(id uint, updates *models.City) (*models.City, error) {
	var city models.City
	if err := s.db.First(&city, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&city).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &city, nil
}

func (s *CityService) DeleteCity(id uint) error {
	return s.db.Delete(&models.City{}, id).Error
}

func (s *CityService) GetCitiesByCountry(country string) ([]models.City, error) {
	var cities []models.City
	err := s.db.Where("country = ? AND status = ?", country, 0).Find(&cities).Error
	return cities, err
}