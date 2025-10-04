package models

import "gorm.io/gorm"

type Insurance struct {
	gorm.Model
	Name        string      `gorm:"type:varchar(100);not null;index" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Coverage    string      `gorm:"type:text" json:"coverage"`
	Phone       string      `gorm:"type:varchar(50)" json:"phone"`
	Email       string      `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Website     string      `gorm:"type:varchar(255)" json:"website"`
	LogoURL     string      `gorm:"type:varchar(255)" json:"logo_url"`
	Cities      []City      `gorm:"many2many:insurance_cities;" json:"cities"` // relation
	Providers   []*Provider `gorm:"many2many:provider_insurances;" json:"providers,omitempty"`
}
