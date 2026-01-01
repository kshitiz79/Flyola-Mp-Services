package models

import "time"

// MealPlan represents meal plan options
type MealPlan struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Code              string    `json:"code" gorm:"type:varchar(10);unique;not null"` // EP, CP, MAP, AP
	Name              string    `json:"name" gorm:"not null"`
	Description       string    `json:"description"`
	IncludesBreakfast bool      `json:"includesBreakfast" gorm:"column:includes_breakfast;default:false"`
	IncludesLunch     bool      `json:"includesLunch" gorm:"column:includes_lunch;default:false"`
	IncludesDinner    bool      `json:"includesDinner" gorm:"column:includes_dinner;default:false"`
	Status            int       `json:"status" gorm:"default:0"` // 0: Active, 1: Inactive
	CreatedAt         time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
