package models

import "time"

// City represents a destination city
type City struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	State       string    `json:"state" gorm:"type:varchar(255)"`
	Country     string    `json:"country" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:text"`
	ImageURL    string    `json:"image_url" gorm:"type:varchar(500)"`
	Status      int       `json:"status"` // 0: Active, 1: Inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}