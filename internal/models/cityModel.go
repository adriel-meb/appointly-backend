package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name       string      `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Insurances []Insurance `gorm:"many2many:insurance_cities;" json:"insurances"`
}
