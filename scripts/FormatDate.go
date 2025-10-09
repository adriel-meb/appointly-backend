package scripts

import "time"

// FormatDate formats a *time.Time into "YYYY-MM-DD" (ISO short date format).
// If the time is nil, it returns an empty string.
func FormatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format("2006-01-02")
	return &formatted
}
