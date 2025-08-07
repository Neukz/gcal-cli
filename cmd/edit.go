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

		event, err := service.Events.Get(calId, eventId).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve event: %v", err)
		}

		updateEventFromFlags(event, cmd.Flags())

		updatedEvent, err := service.Events.Update(calId, eventId, event).Do()
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

	tz, _ := flags.GetString("tz")
	start, _ := flags.GetString("start")
	if start != "" {
		startRFC, err := toRFC3339(start, tz)
		if err != nil {
			log.Fatalf("Invalid --start: %v", err)
		}

		event.Start.DateTime = startRFC
		event.Start.TimeZone = tz
	}

	end, _ := flags.GetString("end")
	if end != "" {
		endRFC, err := toRFC3339(end, tz)
		if err != nil {
			log.Fatalf("Invalid --end: %v", err)
		}

		event.End.DateTime = endRFC
		event.End.TimeZone = tz
	}

	if tz != "" && start == "" && end == "" {
		fmt.Println("--tz provided without --start and/or --end. Time zone will not be updated.")
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
	editCmd.Flags().String("start", "", "start time, e.g. 2025-07-12 13:00")
	editCmd.Flags().String("end", "", "end time, e.g. 2025-07-12 14:00")
	editCmd.Flags().String("desc", "", "event description")
	editCmd.Flags().String("loc", "", "event location")
	editCmd.Flags().String("tz", "", "time zone (IANA name, e.g. Europe/Warsaw) — this only applies if '--start' and/or '--end' are specified, defaults to system time zone")
	editCmd.Flags().String("cal", "", "name of the calendar to edit the event from, defaults to primary")

	editCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(editCmd)
}
