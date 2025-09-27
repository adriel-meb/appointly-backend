package scripts

import (
	"fmt"
	"github.com/adriel-meb/appointly-backend/internal/models"
)

func SendNotifications(notificationType models.NotificationType, message string) error {

	// Send notification asynchronously
	switch notificationType {
	case models.Email:
		fmt.Println(".......Email NOTIFICATION SENT......", message)
	case models.SMS:
		fmt.Println(".......SMS NOTIFICATION SENT........")
	case models.Push:
		fmt.Println(".......Push NOTIFICATION SENT........")
	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)

	}
	return nil
}
