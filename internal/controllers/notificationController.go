package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/adriel-meb/appointly-backend/scripts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateNotification(c *gin.Context, userID uint, message string, notificationType models.NotificationType) error {
	notification := models.Notification{
		UserID:           userID,
		Message:          message,
		NotificationType: notificationType,
	}

	// Save notification in DB
	if err := db.DB.Create(&notification).Error; err != nil {
		return err
	}

	// Try sending notification
	if err := scripts.SendNotifications(notificationType, message); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to send notification",
			Error:   err.Error(),
		})
		return err
	}

	return nil
}
