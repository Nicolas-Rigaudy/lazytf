package aws

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// SSOSession represents an AWS SSO session configuration
type SSOSession struct {
	Name     string
	StartURL string
	Region   string
	Scopes   string
}

// DiscoverSSOSessions parses ~/.aws/config and returns all [sso-session ...] entries
func DiscoverSSOSessions() ([]*SSOSession, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open AWS config file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sessions []*SSOSession
	var currentSession *SSOSession

	for scanner.Scan() {
		// Skip empty lines and comments
		if scanner.Text() == "" || strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[sso-session ") && strings.HasSuffix(line, "]") {
			// Start a new session, and save its name to the slice
			sessionName := strings.TrimSuffix(strings.TrimPrefix(line, "[sso-session "), "]")
			currentSession = &SSOSession{Name: sessionName}
			sessions = append(sessions, currentSession)

		} else if currentSession != nil && strings.Contains(line, "=") {
			// Parse key-value pairs within the current session
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "sso_start_url":
				currentSession.StartURL = value
			case "sso_region":
				currentSession.Region = value
			case "sso_registration_scopes":
				currentSession.Scopes = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading AWS config file: %v", err)
	}

	return sessions, nil
}

// IsLoggedIn checks if AWS SSO credentials are available and valid
// Returns true if credentials work, false otherwise
func IsLoggedIn() bool {
	// Run a quick AWS command with a 2-second timeout to check credentials
	// This is faster than the default timeout and validates the credentials actually work
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity", "--output", "json")
	err := cmd.Run()

	// If the command succeeded, credentials are valid
	return err == nil
}
