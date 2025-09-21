package models

import "gorm.io/gorm"

type NotificationType string

const (
	Email NotificationType = "email"
	SMS   NotificationType = "sms"
	Push  NotificationType = "push"
)

type Notification struct {
	gorm.Model
	UserID           uint             `gorm:"not null"`
	Message          string           `gorm:"type:text" json:"message"`
	NotificationType NotificationType `json:"notification_type"`
}
