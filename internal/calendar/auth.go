package calendar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

const (
	tokenFile     = "token.json"
	serverAddress = ":8080"
	redirectURL   = "http://localhost" + serverAddress
)

// Attempts to retrieve the token and checks validity
func LoadToken() (*oauth2.Token, error) {
	token, err := readTokenFromFile(tokenFile)
	if err != nil {
		return nil, errors.New("token not found")
	}

	if !token.Valid() {
		return nil, errors.New("token expired or invalid")
	}

	return token, nil
}

// Saves the token to file
func SaveToken(token *oauth2.Token) error {
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

// Returns requested token
func StartAuthFlow(config *oauth2.Config) error {
	config.RedirectURL = redirectURL

	// Generate the auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Opening browser to:", authURL)
	go openBrowser(authURL)

	codeCh := make(chan string)
	server := &http.Server{Addr: serverAddress}

	// Listen for redirect on successful authentication
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization failed: no code in request.", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Authorization successful. You can close this window.")
		codeCh <- code

		// Shutdown server after handling
		go server.Shutdown(context.Background())
	})

	serverErrCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	// Wait for code or server error
	select {
	case code := <-codeCh:
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			return err
		}

		return SaveToken(token)
	case err := <-serverErrCh:
		if err != nil {
			return fmt.Errorf("HTTP server error: %v", err)
		}

		return errors.New("server closed unexpectedly")
	}
}

// Retrieves the token and sends a POST to revocation endpoint to invalidate it
func Logout() error {
	token, err := LoadToken()
	if err != nil {
		return err
	}

	// Delete token file
	if err := os.Remove(tokenFile); err != nil {
		return err
	}

	const url = "https://accounts.google.com/o/oauth2/revoke?token="
	const contentType = "application/x-www-form-urlencoded"

	// Helper to revoke a single token
	revokeToken := func(tok string) error {
		res, err := http.Post(url+tok, contentType, nil)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to revoke token, HTTP status: %s", res.Status)
		}

		return nil
	}

	// Try refresh token first
	if err := revokeToken(token.RefreshToken); err == nil {
		return nil
	}

	// If revoking refresh token failed, try access token
	if err := revokeToken(token.AccessToken); err == nil {
		return nil
	}

	// Both failed
	return fmt.Errorf("failed to revoke both refresh and access tokens")
}

// Retrieves the token from file
func readTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Handles opening the browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // Linux
		cmd = "xdg-open"
	}

	if cmd == "rundll32" {
		return exec.Command(cmd, args...).Start()
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
