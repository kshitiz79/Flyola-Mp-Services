package services

import (
	"flyola-services/internal/models"

	"gorm.io/gorm"
)

type MealPlanService struct {
	db *gorm.DB
}

func NewMealPlanService(db *gorm.DB) *MealPlanService {
	return &MealPlanService{db: db}
}

func (s *MealPlanService) GetAllMealPlans() ([]models.MealPlan, error) {
	var mealPlans []models.MealPlan
	err := s.db.Where("status = ?", 0).Find(&mealPlans).Error
	return mealPlans, err
}

func (s *MealPlanService) GetMealPlanByID(id uint) (*models.MealPlan, error) {
	var mealPlan models.MealPlan
	err := s.db.First(&mealPlan, id).Error
	if err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

func (s *MealPlanService) CreateMealPlan(mealPlan *models.MealPlan) error {
	return s.db.Create(mealPlan).Error
}

func (s *MealPlanService) UpdateMealPlan(id uint, updates *models.MealPlan) (*models.MealPlan, error) {
	var mealPlan models.MealPlan
	if err := s.db.First(&mealPlan, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&mealPlan).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

func (s *MealPlanService) DeleteMealPlan(id uint) error {
	return s.db.Delete(&models.MealPlan{}, id).Error
}

func (s *MealPlanService) GetMealPlansByType(planType string) ([]models.MealPlan, error) {
	var mealPlans []models.MealPlan
	err := s.db.Where("type = ? AND status = ?", planType, 0).Find(&mealPlans).Error
	return mealPlans, err
}

func (s *MealPlanService) GetActiveMealPlans() ([]models.MealPlan, error) {
	var mealPlans []models.MealPlan
	err := s.db.Where("status = ? AND is_active = ?", 0, true).Find(&mealPlans).Error
	return mealPlans, err
}