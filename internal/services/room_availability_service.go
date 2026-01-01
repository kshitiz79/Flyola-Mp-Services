package services

import (
	"flyola-services/internal/models"
	"time"

	"gorm.io/gorm"
)

type RoomAvailabilityService struct {
	db *gorm.DB
}

func NewRoomAvailabilityService(db *gorm.DB) *RoomAvailabilityService {
	return &RoomAvailabilityService{db: db}
}

func (s *RoomAvailabilityService) GetAllRoomAvailability(roomID, date string) ([]models.RoomAvailability, error) {
	var availabilities []models.RoomAvailability
	
	query := s.db.Model(&models.RoomAvailability{})
	
	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	
	if date != "" {
		query = query.Where("date = ?", date)
	}
	
	if err := query.Find(&availabilities).Error; err != nil {
		return nil, err
	}

	// Load room and hotel data for each availability record
	for i := range availabilities {
		var room models.Room
		if err := s.db.Preload("Hotel").Preload("RoomCategory").First(&room, availabilities[i].RoomID).Error; err == nil {
			availabilities[i].Room = room
		}
	}

	return availabilities, nil
}

func (s *RoomAvailabilityService) GetRoomAvailabilityByID(id uint) (*models.RoomAvailability, error) {
	var availability models.RoomAvailability
	if err := s.db.First(&availability, id).Error; err != nil {
		return nil, err
	}

	// Load room with hotel and category data
	var room models.Room
	if err := s.db.Preload("Hotel").Preload("RoomCategory").First(&room, availability.RoomID).Error; err == nil {
		availability.Room = room
	}

	return &availability, nil
}

func (s *RoomAvailabilityService) CreateRoomAvailability(availability *models.RoomAvailability) error {
	return s.db.Create(availability).Error
}

func (s *RoomAvailabilityService) UpdateRoomAvailability(id uint, updates *models.RoomAvailability) (*models.RoomAvailability, error) {
	var availability models.RoomAvailability
	if err := s.db.First(&availability, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&availability).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with room and hotel data
	var room models.Room
	if err := s.db.Preload("Hotel").Preload("RoomCategory").First(&room, availability.RoomID).Error; err == nil {
		availability.Room = room
	}

	return &availability, nil
}

func (s *RoomAvailabilityService) DeleteRoomAvailability(id uint) error {
	return s.db.Delete(&models.RoomAvailability{}, id).Error
}

func (s *RoomAvailabilityService) GetAvailabilityByDateRange(roomID uint, startDate, endDate time.Time) ([]models.RoomAvailability, error) {
	var availabilities []models.RoomAvailability
	err := s.db.Where("room_id = ? AND date BETWEEN ? AND ?", roomID, startDate, endDate).Find(&availabilities).Error
	return availabilities, err
}

func (s *RoomAvailabilityService) GetAvailableRoomsByDate(hotelID uint, date time.Time) ([]models.Room, error) {
	var rooms []models.Room
	
	// Get all rooms for the hotel
	if err := s.db.Preload("Hotel").Preload("RoomCategory").Where("hotel_id = ?", hotelID).Find(&rooms).Error; err != nil {
		return nil, err
	}

	// Filter rooms that are available on the given date
	var availableRooms []models.Room
	for _, room := range rooms {
		var availability models.RoomAvailability
		err := s.db.Where("room_id = ? AND date = ? AND is_available = ?", room.ID, date, true).First(&availability).Error
		if err == nil {
			availableRooms = append(availableRooms, room)
		}
	}

	return availableRooms, nil
}

func (s *RoomAvailabilityService) BulkUpdateAvailability(roomID uint, startDate, endDate time.Time, isAvailable bool, price *float64) error {
	updates := map[string]interface{}{
		"is_available": isAvailable,
	}
	
	if price != nil {
		updates["price"] = *price
	}

	return s.db.Model(&models.RoomAvailability{}).
		Where("room_id = ? AND date BETWEEN ? AND ?", roomID, startDate, endDate).
		Updates(updates).Error
}