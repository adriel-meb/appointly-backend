package scripts

import (
	"fmt"
	"github.com/adriel-meb/appointly-backend/internal/models"
)

// GetUpcomingSlotsForProvider returns all available slots for a provider.
// It combines recurring weekly availabilities and one-time date-specific availabilities.
func GetUpcomingSlotsForProvider(providerID uint, daysAhead int) ([]models.AvailabilitySlot, error) {
	fmt.Println("nothing", daysAhead)
	return nil, nil
}
