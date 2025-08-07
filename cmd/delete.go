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
		eventID, _ := cmd.Flags().GetString("id")
		calName, _ := cmd.Flags().GetString("cal")

		if eventID == "" {
			log.Fatalf("Missing required flags: --id.")
		}

		calID, err := calendar.ResolveCalendarID(service, calName)
		if err != nil {
			log.Fatalf("Failed to resolve calendar ID: %v", err)
		}

		err = service.Events.Delete(calID, eventID).Do()
		if err != nil {
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
