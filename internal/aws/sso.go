package aws

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
