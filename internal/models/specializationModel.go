package models

type Specialization struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}
