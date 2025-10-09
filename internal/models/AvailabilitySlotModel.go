package models

import (
	"time"

	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	ID uint `gorm:"primaryKey"` // primary key

	AvailabilityID uint         `gorm:"not null;index"`
	Availability   Availability `gorm:"foreignKey:AvailabilityID;references:ID;constraint:OnDelete:CASCADE"`

	StartTime string `gorm:"type:varchar(5);not null"` // "09:00"
	EndTime   string `gorm:"type:varchar(5);not null"` // "09:30"

	IsBooked bool       `gorm:"default:false"`
	BookedAt *time.Time `gorm:"default:null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
