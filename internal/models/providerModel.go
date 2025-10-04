package models

import (
	"gorm.io/gorm"
)

type Provider struct {
	gorm.Model

	// Basic Info
	UserID           uint            `json:"user_id"`
	User             *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SpecializationID uint            `json:"specialization_id"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization,omitempty"`
	Bio              string          `json:"bio"`

	// Media / Image
	ImageURL string `json:"image_url"` // URL to profile picture
	// Optional: you can add multiple images via a separate Media table

	// Insurances (many-to-many)
	Insurances []*Insurance `gorm:"many2many:provider_insurances;" json:"insurances,omitempty"`

	// Availabilities (one-to-many)
	Availabilities []Availability `json:"availabilities,omitempty"`

	// Computed field (not stored in DB)
	NextAvailable string `gorm:"-" json:"next_available,omitempty"`

	//city
	CityID uint  `gorm:"not null"`
	City   *City `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"city,omitempty"`

	// Optional: Ratings, Reviews, Price, etc.
	Rating      float32 `json:"rating,omitempty"`
	ReviewCount int     `json:"review_count,omitempty"`
	Price       string  `json:"price,omitempty"`
	Address     string  `json:"address,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lng         float64 `json:"lng,omitempty"`
	Distance    string  `json:"distance"`
}
