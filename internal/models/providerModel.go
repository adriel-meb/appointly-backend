package models

import (
	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model
	UserID           uint            `json:"user_id"`
	User             *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SpecializationID uint            `json:"specialization_id"` // FK
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization,omitempty"`
	Bio              string          `json:"bio"`
}
