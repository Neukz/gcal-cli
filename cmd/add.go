package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
	cal "google.golang.org/api/calendar/v3"
)

// Flags
var (
	title   string
	start   string
	end     string
	desc    string
	loc     string
	tz      string
	calName string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an event to Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		service := calendar.GetService()

		// Required flags
		if title == "" || start == "" || end == "" {
			log.Fatalf("Missing required flags: --title, --start, and --end.")
		}

		startRFC, err := toRFC3339(start, tz)
		if err != nil {
			log.Fatalf("Invalid --start: %v", err)
		}

		endRFC, err := toRFC3339(end, tz)
		if err != nil {
			log.Fatalf("Invalid --end: %v", err)
		}

		calID := "primary" // Default calendar
		if calName != "" {
			id, err := calendar.ResolveCalendarID(service, calName)
			if err != nil {
				log.Fatalf("Failed to resolve calendar ID: %v", err)
			}

			calID = id
		}

		event := &cal.Event{
			Summary:     title,
			Description: desc,
			Location:    loc,
			Start: &cal.EventDateTime{
				DateTime: startRFC,
				TimeZone: tz,
			},
			End: &cal.EventDateTime{
				DateTime: endRFC,
				TimeZone: tz,
			},
		}

		createdEvent, err := service.Events.Insert(calID, event).Do()
		if err != nil {
			log.Fatalf("Unable to create event: %v", err)
		}

		fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)
	},
}

// Parses the datetime in "YYYY-MM-DD HH:MM" format and returns it in RFC3339
func toRFC3339(datetime, tz string) (string, error) {
	const layout = "2006-01-02 15:04"

	loc := time.Now().Location() // Default to system time zone
	var err error

	if tz != "" {
		loc, err = time.LoadLocation(tz)
		if err != nil {
			return "", fmt.Errorf("invalid timezone: %w", err)
		}
	}

	t, err := time.ParseInLocation(layout, datetime, loc)
	if err != nil {
		return "", fmt.Errorf("invalid datetime format: %w", err)
	}

	return t.Format(time.RFC3339), nil
}

func init() {
	addCmd.Flags().StringVarP(&title, "title", "t", "", "event title (required)")
	addCmd.Flags().StringVarP(&start, "start", "s", "", "start time, e.g. 2025-07-12 13:00 (required)")
	addCmd.Flags().StringVarP(&end, "end", "e", "", "end time, e.g. 2025-07-12 14:00 (required)")
	addCmd.Flags().StringVarP(&desc, "desc", "d", "", "event description")
	addCmd.Flags().StringVarP(&loc, "loc", "l", "", "event location")
	addCmd.Flags().StringVarP(&tz, "tz", "z", "", "time zone (IANA name, e.g. Europe/Warsaw), defaults to system time zone")
	addCmd.Flags().StringVarP(&calName, "cal", "c", "", "name of the calendar to insert the event into, defaults to primary")

	rootCmd.MarkFlagRequired("title")
	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("end")

	rootCmd.AddCommand(addCmd)
}
