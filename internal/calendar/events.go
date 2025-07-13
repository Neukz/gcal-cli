package calendar

import "google.golang.org/api/calendar/v3"

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
