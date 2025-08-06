package calendar

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const CredentialsFile = "credentials.json"

func GetService() *calendar.Service {
	b, err := os.ReadFile(CredentialsFile)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	token, err := LoadToken()
	if err != nil {
		log.Fatalf("No valid token found. Run `gcal login` first.")
	}

	client := config.Client(context.Background(), token)

	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create calendar service: %v", err)
	}

	return service
}
