package calendar

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

// Returns an event created from arguments
func NewEvent(title, desc, loc, tz, startRFC, endRFC string) *calendar.Event {
	return &calendar.Event{
		Summary:     title,
		Description: desc,
		Location:    loc,
		Start: &calendar.EventDateTime{
			DateTime: startRFC,
			TimeZone: tz,
		},
		End: &calendar.EventDateTime{
			DateTime: endRFC,
			TimeZone: tz,
		},
	}
}

// Fetches events between tMin and tMax for the given calendar ID
func GetEvents(service *calendar.Service, calId string, tMin, tMax time.Time, maxResults int) ([]*calendar.Event, error) {
	// Prepare request
	req := service.Events.List(calId).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(tMin.Format(time.RFC3339)).
		TimeMax(tMax.Format(time.RFC3339)).
		OrderBy("startTime")

	// Limit results
	if maxResults > 0 {
		req = req.MaxResults(int64(maxResults))
	}

	events, err := req.Do()
	if err != nil {
		return nil, err
	}

	return events.Items, nil
}
