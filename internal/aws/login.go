package aws

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
	tea "github.com/charmbracelet/bubbletea"
)

// RunSSOLogin executes `aws sso login --sso-session <sessionName>`
// This command opens a browser for authentication
func RunSSOLogin(session *SSOSession) tea.Cmd {
	args := []string{"sso", "login", "--sso-session", session.Name}
	result := executor.ExecuteStreaming("aws", args, "", false)
	return result.Cmd
}
