package controllers

import (
	"fmt"
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
)

func CreateNotification(userID uint, message string, NotificationType models.NotificationType) error {

	notification := models.Notification{
		UserID:           userID,
		Message:          message,
		NotificationType: NotificationType,
	}

	if err := db.DB.Create(&notification).Error; err != nil {
		return err
	}

	// Send notification asynchronously
	switch NotificationType {
	case models.Email:
		fmt.Println(".......Email NOTIFICATION SENT........")
	case models.SMS:
		fmt.Println(".......SMS NOTIFICATION SENT........")
	case models.Push:
		fmt.Println(".......Push NOTIFICATION SENT........")

	}

	return nil
}
