package scripts

import (
	"fmt"
	"time"
)

// GenerateTimeSlots divides a time range into equal slots (e.g. 30min).
// GenerateTimeSlots divides a time range into equal slots (e.g., 30min).
func GenerateTimeSlots(startStr, endStr string, slotMinutes int, includeStart, includeEnd bool) ([]string, error) {
	if slotMinutes <= 0 {
		return nil, fmt.Errorf("slotMinutes must be positive, got %d", slotMinutes)
	}

	layout := "15:04"
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %v", err)
	}
	end, err := time.Parse(layout, endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %v", err)
	}

	if !end.After(start) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	var slots []string
	current := start

	// Include start time if requested
	if includeStart {
		slots = append(slots, current.Format(layout))
	}

	// Generate slots
	for {
		current = current.Add(time.Duration(slotMinutes) * time.Minute)

		// Check if we should include this slot
		if current.After(end) {
			break
		}
		if current.Equal(end) && !includeEnd {
			break
		}

		slots = append(slots, current.Format(layout))
	}

	return slots, nil
}
