package models

import (
	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	ID               uint           `gorm:"primaryKey"`
	UserID           uint           `gorm:"not null"`
	SpecializationID uint           `gorm:"not null"`
	Specialization   Specialization `gorm:"foreignKey:SpecializationID"`
	Bio              string
}
