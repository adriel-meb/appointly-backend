package scripts

import (
	"fmt"
	"time"

	"github.com/adriel-meb/appointly-backend/internal/models"
)

// ValidateAvailabilityInput validates availability fields and overlapping slots
func ValidateAvailabilityInput(input models.Availability, existing []models.Availability) error {
	// 1️⃣ Recurring vs one-time checks
	if input.IsRecurring {
		if input.DayOfWeek == nil {
			return fmt.Errorf("DayOfWeek is required for recurring availability")
		}
		if input.Date != nil {
			return fmt.Errorf("Date must not be provided for recurring availability")
		}
	} else {
		if input.Date == nil {
			return fmt.Errorf("Date is required for one-time availability")
		}
		if input.DayOfWeek != nil {
			return fmt.Errorf("DayOfWeek must not be provided for one-time availability")
		}
	}

	// 2️⃣ Validate StartTime < EndTime
	start, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return fmt.Errorf("StartTime invalid: %v", err)
	}
	end, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return fmt.Errorf("EndTime invalid: %v", err)
	}
	if !end.After(start) {
		return fmt.Errorf("EndTime must be after StartTime")
	}

	// 3️⃣ Check for overlapping slots
	for _, e := range existing {
		existingStart, _ := time.Parse("15:04", e.StartTime)
		existingEnd, _ := time.Parse("15:04", e.EndTime)

		if start.Before(existingEnd) && end.After(existingStart) {
			return fmt.Errorf("Availability overlaps with an existing slot")
		}
	}

	return nil
}
