package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Base command
var rootCmd = &cobra.Command{
	Use:   "gcal",
	Short: "CLI for Google Calendar",
}

// Runs the base command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
