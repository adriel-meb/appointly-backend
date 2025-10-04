package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusBooking string

// Booking statuses

const (
	Pending   StatusBooking = "pending"
	Confirmed StatusBooking = "confirmed"
	Completed StatusBooking = "completed"
	Cancelled StatusBooking = "cancelled"
)

type Booking struct {
	gorm.Model

	PatientID      uint  `gorm:"not null;index" json:"patient_id"`  // search bookings by patient
	ProviderID     uint  `gorm:"not null;index" json:"provider_id"` // search bookings by provider
	ServiceID      uint  `gorm:"not null;index" json:"service_id"`  // search bookings by service
	AvailabilityID *uint `gorm:"index" json:"availability_id"`

	// Timing
	StartTime time.Time `gorm:"not null;index" json:"start_time"` // useful for availability queries
	EndTime   time.Time `gorm:"index" json:"end_time"`

	// Status & details
	Status StatusBooking `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Notes  string        `gorm:"type:text" json:"notes"`

	// Optional payment tracking
	PaymentStatus string  `gorm:"type:varchar(20);default:'unpaid';index" json:"payment_status"`
	Amount        float64 `json:"amount"`
}
