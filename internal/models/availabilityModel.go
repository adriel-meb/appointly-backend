package models

import (
	"gorm.io/gorm"
	"time"
)

// DayOfWeekEnum represents allowed days for recurring availability
type DayOfWeekEnum string

const (
	Monday    DayOfWeekEnum = "MONDAY"
	Tuesday   DayOfWeekEnum = "TUESDAY"
	Wednesday DayOfWeekEnum = "WEDNESDAY"
	Thursday  DayOfWeekEnum = "THURSDAY"
	Friday    DayOfWeekEnum = "FRIDAY"
	Saturday  DayOfWeekEnum = "SATURDAY"
	Sunday    DayOfWeekEnum = "SUNDAY"
)

// Availability represents a provider's available time slot
type Availability struct {
	gorm.Model

	ProviderID uint     `gorm:"not null" json:"provider_id"`
	Provider   Provider `gorm:"foreignKey:ProviderID"`

	// For recurring weekly slots
	DayOfWeek   *DayOfWeekEnum `gorm:"type:varchar(10)" json:"day_of_week"`
	IsRecurring bool           `gorm:"default:false" json:"is_recurring"`

	// For specific one-time slots
	Date *time.Time `json:"date"`

	// Time range
	StartTime string `gorm:"type:varchar(5);not null" json:"start_time"` // "09:00"
	EndTime   string `gorm:"type:varchar(5);not null" json:"end_time"`   // "17:00"
}
