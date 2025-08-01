package cmd

import (
	"fmt"
	"log"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cal "google.golang.org/api/calendar/v3"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an event",
	Run: func(cmd *cobra.Command, args []string) {
		service := calendar.GetService()

		// Flags
		eventID, _ := cmd.Flags().GetString("id")
		calName, _ := cmd.Flags().GetString("cal")

		if eventID == "" {
			log.Fatalf("Missing required flags: --id.")
		}

		calID, err := calendar.ResolveCalendarID(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		event, err := service.Events.Get(calID, eventID).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve event: %v", err)
		}

		updateEventFromFlags(event, cmd.Flags())

		updatedEvent, err := service.Events.Update(calID, eventID, event).Do()
		if err != nil {
			log.Fatalf("Unable to edit event: %v", err)
		}

		fmt.Printf("Event updated: %s\n", updatedEvent.HtmlLink)
	},
}

// Updates event information based on flags
func updateEventFromFlags(event *cal.Event, flags *pflag.FlagSet) {
	title, _ := flags.GetString("title")
	if title != "" {
		event.Summary = title
	}

	desc, _ := flags.GetString("desc")
	if desc != "" {
		event.Description = desc
	}

	loc, _ := flags.GetString("loc")
	if loc != "" {
		event.Location = loc
	}
}

func init() {
	editCmd.Flags().String("id", "", "event ID (required)")
	editCmd.Flags().String("title", "", "event title")
	editCmd.Flags().String("desc", "", "event description")
	editCmd.Flags().String("loc", "", "event location")
	editCmd.Flags().String("cal", "", "name of the calendar to edit the event from, defaults to primary")

	editCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(editCmd)
}
