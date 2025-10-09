package models

import (
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

// Availability represents a provider's available time range
type Availability struct {
	ID uint `gorm:"primaryKey"` // primary key

	ProviderID uint     `gorm:"not null;index"`
	Provider   Provider `gorm:"foreignKey:ProviderID"`

	// Recurring weekly slot
	DayOfWeek   *DayOfWeekEnum `gorm:"type:varchar(10)" json:"day_of_week"`
	IsRecurring bool           `gorm:"default:false" json:"is_recurring"`

	// One-time slot
	Date *time.Time `json:"date"`

	// Time range
	StartTime string `gorm:"type:varchar(5);not null" json:"start_time"` // "09:00"
	EndTime   string `gorm:"type:varchar(5);not null" json:"end_time"`   // "17:00"

	Slots []AvailabilitySlot `gorm:"foreignKey:AvailabilityID;constraint:OnDelete:CASCADE"` // cascade deletion
}
