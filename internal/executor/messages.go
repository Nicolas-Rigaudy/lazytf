package executor

import tea "github.com/charmbracelet/bubbletea"

// CommandOutputMsg represents a single line of output from a running command
type CommandOutputMsg struct {
	Line       string
	IsErr      bool    // true for stderr, false for stdout
	ListenNext tea.Cmd // Command to listen for next message
}

// CommandCompletedMsg is sent when a command finishes successfully
type CommandCompletedMsg struct {
	Command  string
	ExitCode int
	Output   string
}

// CommandErrorMsg is sent when a command fails
type CommandErrorMsg struct {
	Command string
	Error   error
	Output  string
}
