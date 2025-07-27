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
		service := calendar.GetService()

		// Flags
		eventID, _ := cmd.Flags().GetString("id")

		if eventID == "" {
			log.Fatalf("Missing required flags: --id.")
		}

		err := service.Events.Delete("primary", eventID).Do()
		if err != nil {
			log.Fatalf("Unable to delete event: %v", err)
		}

		fmt.Println("Event deleted.")
	},
}

func init() {
	deleteCmd.Flags().String("id", "", "event ID (required)")

	deleteCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(deleteCmd)
}
