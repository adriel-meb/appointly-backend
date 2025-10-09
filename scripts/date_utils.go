package scripts

import (
	"errors"
	"time"
)

// ParseDateFlexible tries multiple date formats safely
func ParseDateFlexible(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, errors.New("empty date string")
	}

	formats := []string{
		"2006-01-02",       // e.g. 2025-10-12
		"02-01-2006",       // e.g. 12-10-2025
		time.RFC3339,       // e.g. 2025-10-12T00:00:00Z
		"2006-01-02T15:04", // e.g. 2025-10-12T09:00
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, errors.New("invalid date format: expected YYYY-MM-DD or DD-MM-YYYY")
}
