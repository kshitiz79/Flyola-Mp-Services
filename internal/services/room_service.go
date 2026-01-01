package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type RoomService struct {
	db *gorm.DB
}

func NewRoomService(db *gorm.DB) *RoomService {
	return &RoomService{db: db}
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.Preload("Hotel").Preload("RoomCategory").Find(&rooms).Error
	return rooms, err
}

func (s *RoomService) GetRoomByID(id uint) (*models.Room, error) {
	var room models.Room
	err := s.db.Preload("Hotel").Preload("RoomCategory").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *RoomService) CreateRoom(room *models.Room) error {
	return s.db.Create(room).Error
}

func (s *RoomService) UpdateRoom(id uint, updates *models.Room) (*models.Room, error) {
	var room models.Room
	if err := s.db.First(&room, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&room).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Hotel").Preload("RoomCategory").First(&room, id).Error; err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *RoomService) DeleteRoom(id uint) error {
	return s.db.Delete(&models.Room{}, id).Error
}

func (s *RoomService) GetRoomsByHotel(hotelID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.Preload("Hotel").Preload("RoomCategory").Where("hotel_id = ?", hotelID).Find(&rooms).Error
	return rooms, err
}

func (s *RoomService) GetRoomsByCategory(categoryID uint) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.Preload("Hotel").Preload("RoomCategory").Where("room_category_id = ?", categoryID).Find(&rooms).Error
	return rooms, err
}

func (s *RoomService) GetAvailableRooms(hotelID uint, status int) ([]models.Room, error) {
	var rooms []models.Room
	err := s.db.Preload("Hotel").Preload("RoomCategory").Where("hotel_id = ? AND status = ?", hotelID, status).Find(&rooms).Error
	return rooms, err
}

func (s *RoomService) UpdateRoomStatus(id uint, status int) error {
	return s.db.Model(&models.Room{}).Where("id = ?", id).Update("status", status).Error
}