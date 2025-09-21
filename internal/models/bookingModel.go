package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusBooking string

const (
	Pending   StatusBooking = "pending"
	Confirmed StatusBooking = "confirmed"
	Completed StatusBooking = "completed"
	Cancelled StatusBooking = "cancelled"
)

type Booking struct {
	gorm.Model

	PatientID      uint  `gorm:"not null" json:"patient_id"`
	ProviderID     uint  `gorm:"not null" json:"provider_id"`
	ServiceID      uint  `gorm:"not null" json:"service_id"`
	AvailabilityID *uint `json:"availability_id"` // optional link to availability

	// Timing
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	// Status & details
	Status StatusBooking `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Notes  string        `gorm:"type:text" json:"notes"`

	// Optional payment tracking
	PaymentStatus string  `gorm:"type:varchar(20);default:'unpaid'" json:"payment_status"`
	Amount        float64 `json:"amount"`
}
