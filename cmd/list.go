package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
	cal "google.golang.org/api/calendar/v3"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Google Calendar events",
	Run: func(cmd *cobra.Command, args []string) {
		service := calendar.GetService()

		// Flags
		tomorrow, _ := cmd.Flags().GetBool("tomorrow")
		daysAhead, _ := cmd.Flags().GetInt("days")
		maxResults, _ := cmd.Flags().GetInt("max")
		calName, _ := cmd.Flags().GetString("cal")

		tNow := time.Now()
		var tMin, tMax time.Time
		var heading string

		// Handle mutally exclusive flags
		switch {
		case tomorrow:
			tMin = time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, tNow.Location()).Add(24 * time.Hour)
			tMax = tMin.Add(24 * time.Hour)
			heading = "Events for tomorrow:"
		case daysAhead > 0:
			tMin = tNow
			tMax = tNow.Add(time.Duration(daysAhead) * 24 * time.Hour)
			heading = fmt.Sprintf("Events for the next %d day(s):", daysAhead)
		default:
			tMin = time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, tNow.Location())
			tMax = tMin.Add(24 * time.Hour)
			heading = "Events for today:"
		}

		calID, err := calendar.ResolveCalendarID(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		events, err := calendar.GetEvents(service, calID, tMin, tMax, maxResults)
		if err != nil {
			log.Fatalf("Unable to retrieve events: %v", err)
		}

		if len(events) == 0 {
			fmt.Println("No events found.")
			return
		}

		fmt.Println(heading)
		for _, event := range events {
			printEvent(event)
		}

		if maxResults > 0 && len(events) == maxResults {
			fmt.Printf("Stopped after first %d entries.\n", maxResults)
		}
	},
}

// Prints formatted event date, time, and title
func printEvent(event *cal.Event) {
	var tReadable string

	if event.Start.DateTime != "" { // Time-specific event
		t, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			tReadable = event.Start.DateTime // Fallback
		} else {
			tReadable = t.Format("Mon, Jan 2 at 15:04")
		}
	} else if event.Start.Date != "" { // All-day event
		t, err := time.Parse("2006-01-02", event.Start.Date)
		if err != nil {
			tReadable = event.Start.Date // Fallback
		} else {
			tReadable = t.Format("Mon, Jan 2") + " (All day)"
		}
	} else {
		tReadable = "Unknown time"
	}

	fmt.Printf("> %s — %s\n", tReadable, event.Summary)
}

func init() {
	listCmd.Flags().Bool("tomorrow", false, "show events for tomorrow")
	listCmd.Flags().Int("days", 0, "number of days ahead to list events")
	listCmd.Flags().Int("max", 0, "maximum number of events to show")
	listCmd.Flags().String("cal", "", "name of the calendar to list events from, defaults to primary")

	rootCmd.AddCommand(listCmd)
}
