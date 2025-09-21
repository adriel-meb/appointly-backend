package models

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Title           string   `gorm:"type:varchar(100);not null"`
	ProviderID      uint     `gorm:"not null"`
	Provider        Provider `gorm:"foreignKey:ProviderID"`
	Description     string   `gorm:"type:text"`
	DurationMinutes uint     `gorm:"not null;check:duration_minutes > 0"`
	Price           float64  `gorm:"not null;check:price > 0"`
}
