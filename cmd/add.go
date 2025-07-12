package cmd

import (
	"fmt"
	"log"
	"time"

	gcal "github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
)

// Falgs
var (
	title string
	start string
	end   string
	desc  string
	loc   string
	tz    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an event to Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
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

		service := gcal.GetService()

		event := &calendar.Event{
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

		createdEvent, err := service.Events.Insert("primary", event).Do()
		if err != nil {
			log.Fatalf("Unable to create event: %v", err)
		}

		fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)
	},
}

func toRFC3339(datetime, tz string) (string, error) {
	const layout = "2006-01-02 15:04"

	// Use system time zone by default
	loc := time.Now().Location()
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
	// Register flags
	addCmd.Flags().StringVarP(&title, "title", "t", "", "event title (required)")
	addCmd.Flags().StringVarP(&start, "start", "s", "", "start time, e.g. 2025-07-12 13:00 (required)")
	addCmd.Flags().StringVarP(&end, "end", "e", "", "end time, e.g. 2025-07-12 14:00 (required)")
	addCmd.Flags().StringVarP(&desc, "desc", "d", "", "event description")
	addCmd.Flags().StringVarP(&loc, "loc", "l", "", "event location")
	addCmd.Flags().StringVarP(&tz, "tz", "z", "", "time zone (IANA name, e.g. Europe/Warsaw), defaults to system time zone")

	rootCmd.MarkFlagRequired("title")
	rootCmd.MarkFlagRequired("start")
	rootCmd.MarkFlagRequired("end")

	rootCmd.AddCommand(addCmd)
}
