package models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	Specialization string         `gorm:"type:varchar(100);not null" json:"specialization" binding:"required"`
	bio            string         `gorm:"type:text" json:"bio,omitempty"`
	UserID         uint           `gorm:"not null;uniqueIndex" json:"user_id" binding:"required"`
	User           User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // soft deletes
}
