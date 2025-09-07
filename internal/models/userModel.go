package models

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	RolePatient  UserRole = "patient"
	RoleProvider UserRole = "provider"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	gorm.Model
	Name           string   `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Email          string   `gorm:"type:varchar(150);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password       string   `gorm:"-" json:"password,omitempty" binding:"required,min=6"` // input only, ignored by DB
	PasswordHash   string   `gorm:"type:text;not null" json:"-"`                          // stored hash, hidden in API
	Role           UserRole `gorm:"type:varchar(20);not null;default:'patient'" json:"role" binding:"omitempty,oneof=patient provider admin"`
	PhoneNumber    *string  `gorm:"type:varchar(20)" json:"phone,omitempty"` // optional
	Specialization string   `json:"specialization,omitempty"`                // Only for provider
	Bio            string   `json:"bio,omitempty"`                           // Only for provider
}
