package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Base command
var rootCmd = &cobra.Command{
	Use:   "gcal",
	Short: "A CLI for Google Calendar",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Add child commands and flags
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
