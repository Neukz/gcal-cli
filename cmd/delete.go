package cmd

import (
	"fmt"
	"log"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an event from Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := calendar.GetService()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Flags
		eventId, _ := cmd.Flags().GetString("id")
		calName, _ := cmd.Flags().GetString("cal")

		if eventId == "" {
			log.Fatalf("Missing required flags: --id.")
		}

		calId, err := calendar.ResolveCalendarId(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		if err := service.Events.Delete(calId, eventId).Do(); err != nil {
			log.Fatalf("Unable to delete event: %v", err)
		}

		fmt.Println("Event deleted.")
	},
}

func init() {
	deleteCmd.Flags().String("id", "", "event ID (required)")
	deleteCmd.Flags().String("cal", "", "name of the calendar to delete the event from, defaults to primary")

	deleteCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(deleteCmd)
}
