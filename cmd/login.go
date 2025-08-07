package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Neukz/gcal-cli/internal/calendar"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := os.ReadFile(calendar.CredentialsFilename)
		if err != nil {
			log.Fatalf("Unable to read credentials file: %v", err)
		}

		config, err := google.ConfigFromJSON(creds, "https://www.googleapis.com/auth/calendar")
		if err != nil {
			log.Fatalf("Unable to parse credentials: %v", err)
		}

		if err := calendar.StartAuthFlow(config); err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}

		fmt.Println("Authentication successful.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
