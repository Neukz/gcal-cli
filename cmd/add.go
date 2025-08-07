package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new event",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := calendar.GetService()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Flags
		title, _ := cmd.Flags().GetString("title")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		desc, _ := cmd.Flags().GetString("desc")
		loc, _ := cmd.Flags().GetString("loc")
		tz, _ := cmd.Flags().GetString("tz")
		calName, _ := cmd.Flags().GetString("cal")

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

		calId, err := calendar.ResolveCalendarId(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		event := calendar.NewEvent(title, desc, loc, tz, startRFC, endRFC)
		createdEvent, err := service.Events.Insert(calId, event).Do()
		if err != nil {
			log.Fatalf("Unable to create event: %v", err)
		}

		fmt.Printf("Event created.\nID: %s\nLink: %s\n", createdEvent.Id, createdEvent.HtmlLink)
	},
}

// Parses the datetime in "YYYY-MM-DD HH:MM" format and returns it in RFC3339
func toRFC3339(datetime, tz string) (string, error) {
	const layout = "2006-01-02 15:04"

	loc := time.Now().Location() // Defaults to system time zone
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
	addCmd.Flags().String("title", "", "event title (required)")
	addCmd.Flags().String("start", "", "start time, e.g. 2025-07-12 13:00 (required)")
	addCmd.Flags().String("end", "", "end time, e.g. 2025-07-12 14:00 (required)")
	addCmd.Flags().String("desc", "", "event description")
	addCmd.Flags().String("loc", "", "event location")
	addCmd.Flags().String("tz", "", "time zone (IANA name, e.g. Europe/Warsaw), defaults to system time zone")
	addCmd.Flags().String("cal", "", "name of the calendar to insert the event into, defaults to primary")

	addCmd.MarkFlagRequired("title")
	addCmd.MarkFlagRequired("start")
	addCmd.MarkFlagRequired("end")

	rootCmd.AddCommand(addCmd)
}
