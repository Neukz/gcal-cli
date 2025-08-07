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
	Short: "List events",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := calendar.GetService()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Flags
		showAll, _ := cmd.Flags().GetBool("all")
		tomorrow, _ := cmd.Flags().GetBool("tomorrow")
		daysAhead, _ := cmd.Flags().GetInt("days")
		maxResults, _ := cmd.Flags().GetInt("max")
		calName, _ := cmd.Flags().GetString("cal")

		tNow := time.Now()
		tMin := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, tNow.Location()) // Include the full day
		var tMax time.Time
		const day = 24 * time.Hour

		// Handle mutally exclusive flags
		switch {
		case tomorrow:
			tMin = tMin.Add(day)
			tMax = tMin.Add(day)
		case daysAhead > 0:
			tMax = tMin.Add(time.Duration(daysAhead+1) * day)
		default:
			tMax = tMin.Add(day)
		}

		calId, err := calendar.ResolveCalendarId(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		events, err := calendar.GetEvents(service, calId, tMin, tMax, maxResults)
		if err != nil {
			log.Fatalf("Unable to retrieve events: %v", err)
		}

		if len(events) == 0 {
			fmt.Println("No events found.")
			return
		}

		for _, event := range events {
			printEvent(event, showAll)
		}

		if maxResults > 0 && len(events) == maxResults {
			fmt.Printf("Stopped after first %d entries.\n", maxResults)
		}
	},
}

// Prints formatted event date, time, and title
func printEvent(event *cal.Event, showAll bool) {
	var startStr string

	if event.Start.DateTime != "" { // Time-specific event
		t, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			startStr = event.Start.DateTime // Fallback
		} else {
			startStr = t.Format("Mon, Jan 2 at 15:04")
		}
	} else if event.Start.Date != "" { // All-day event
		t, err := time.Parse("2006-01-02", event.Start.Date)
		if err != nil {
			startStr = event.Start.Date // Fallback
		} else {
			startStr = t.Format("Mon, Jan 2") + " (All day)"
		}
	} else {
		startStr = "Unknown time"
	}

	fmt.Printf("> %s — %s\n", startStr, event.Summary)

	if showAll {
		fmt.Printf("\tID: %s\n", event.Id)
		if event.Description != "" {
			fmt.Printf("\tDescription: %s\n", event.Description)
		}
		if event.Location != "" {
			fmt.Printf("\tLocation: %s\n", event.Location)
		}
		if event.HtmlLink != "" {
			fmt.Printf("\tLink: %s\n", event.HtmlLink)
		}
	}
}

func init() {
	listCmd.Flags().Bool("all", false, "display event details")
	listCmd.Flags().Bool("tomorrow", false, "show events for tomorrow")
	listCmd.Flags().Int("days", 0, "number of days ahead to list events")
	listCmd.Flags().Int("max", 0, "maximum number of events to show")
	listCmd.Flags().String("cal", "", "name of the calendar to list events from, defaults to primary")

	rootCmd.AddCommand(listCmd)
}
