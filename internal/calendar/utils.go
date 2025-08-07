package calendar

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

// Finds the calendar ID by its human-readable name
func ResolveCalendarId(service *calendar.Service, name string) (string, error) {
	// Default calendar
	if name == "" {
		return "primary", nil
	}

	cals, err := service.CalendarList.List().Do()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve calendar list: %v", err)
	}

	for _, cal := range cals.Items {
		if cal.Summary == name {
			return cal.Id, nil
		}
	}

	return "", fmt.Errorf("calendar with name %q not found", name)
}
