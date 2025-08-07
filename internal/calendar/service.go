package calendar

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const CredentialsFile = "credentials.json"

func GetService() (*calendar.Service, error) {
	credentials, err := os.ReadFile(CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, calendar.CalendarScope)
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
