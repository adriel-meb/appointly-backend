package scripts

import (
	"fmt"
	"github.com/adriel-meb/appointly-backend/internal/models"
)

// ValidateAvailabilityInput checks the consistency of an availability input.
// Rules:
//  1. If IsRecurring is true → DayOfWeek must be set, Date must be nil.
//  2. If IsRecurring is false → Date must be set, DayOfWeek must be nil.
func ValidateAvailabilityInput(input models.Availability) error {

	// Case 1: recurring availability
	if input.IsRecurring {
		// Ensure a DayOfWeek is provided (e.g., "Monday")
		if input.DayOfWeek == nil {
			return fmt.Errorf("DayOfWeek is required for recurring availability")
		}
		// Ensure a Date is NOT provided for recurring slots
		if input.Date != nil {
			return fmt.Errorf("Date must not be provided for recurring availability")
		}
	} else {
		// Case 2: one-time availability
		// Ensure a specific Date is provided
		if input.Date == nil {
			return fmt.Errorf("Date is required for one-time availability")
		}
		// Ensure DayOfWeek is NOT provided for one-time slots
		if input.DayOfWeek != nil {
			return fmt.Errorf("DayOfWeek must not be provided for one-time availability")
		}
	}

	// If all checks pass, return no error
	return nil
}
