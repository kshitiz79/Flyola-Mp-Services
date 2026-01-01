package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) GetAllReviews() ([]models.HotelReview, error) {
	var reviews []models.HotelReview
	err := s.db.Preload("Hotel").Preload("Booking").Where("status = ?", 0).Find(&reviews).Error
	return reviews, err
}

func (s *ReviewService) GetReviewByID(id uint) (*models.HotelReview, error) {
	var review models.HotelReview
	err := s.db.Preload("Hotel").Preload("Booking").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (s *ReviewService) CreateReview(review *models.HotelReview) error {
	return s.db.Create(review).Error
}

func (s *ReviewService) UpdateReview(id uint, updates *models.HotelReview) (*models.HotelReview, error) {
	var review models.HotelReview
	if err := s.db.First(&review, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&review).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload with relationships
	if err := s.db.Preload("Hotel").Preload("Booking").First(&review, id).Error; err != nil {
		return nil, err
	}

	return &review, nil
}

func (s *ReviewService) DeleteReview(id uint) error {
	return s.db.Delete(&models.HotelReview{}, id).Error
}

func (s *ReviewService) GetHotelReviews(hotelID uint) ([]models.HotelReview, error) {
	var reviews []models.HotelReview
	err := s.db.Preload("Hotel").Preload("Booking").
		Where("hotel_id = ? AND status = ?", hotelID, 0).Find(&reviews).Error
	return reviews, err
}

func (s *ReviewService) UpdateReviewStatus(id uint, status int) (*models.HotelReview, error) {
	var review models.HotelReview
	if err := s.db.First(&review, id).Error; err != nil {
		return nil, err
	}

	review.Status = status
	if err := s.db.Save(&review).Error; err != nil {
		return nil, err
	}

	return &review, nil
}

func (s *ReviewService) GetReviewsByRating(rating int) ([]models.HotelReview, error) {
	var reviews []models.HotelReview
	err := s.db.Preload("Hotel").Preload("Booking").
		Where("rating = ? AND status = ?", rating, 0).Find(&reviews).Error
	return reviews, err
}

func (s *ReviewService) GetVerifiedReviews() ([]models.HotelReview, error) {
	var reviews []models.HotelReview
	err := s.db.Preload("Hotel").Preload("Booking").
		Where("is_verified = ? AND status = ?", 1, 0).Find(&reviews).Error
	return reviews, err
}