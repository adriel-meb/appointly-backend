package models

import (
	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	Specialization string         `gorm:"type:varchar(100);not null" json:"specialization" binding:"required"`
	Bio            string         `gorm:"type:text" json:"bio,omitempty"`
	UserID         uint           `gorm:"not null;uniqueIndex" json:"user_id" binding:"required"`
	User           User           `gorm:"foreignKey:UserID" json:"user"`
}
