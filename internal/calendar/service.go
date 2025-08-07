package calendar

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const CredentialsFilename = "credentials.json"

// Returns the calendar service after successful authentication
func GetService() (*calendar.Service, error) {
	creds, err := os.ReadFile(CredentialsFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(creds, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	token, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the token: %v", err)
	}

	client := config.Client(context.Background(), token)

	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %v", err)
	}

	return service, nil
}
