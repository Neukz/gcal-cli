package cmd

import (
	"fmt"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		if err := calendar.Logout(); err != nil {
			fmt.Printf("Already logged out: %v", err)
			return
		}

		fmt.Println("Logged out successfully.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
