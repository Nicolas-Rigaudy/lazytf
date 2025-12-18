package terraform

import tea "github.com/charmbracelet/bubbletea"

type CommandOutputMsg struct {
	Line       string
	IsErr      bool // true for stderr
	ListenNext tea.Cmd // Command to listen for next message
}

type CommandCompletedMsg struct {
	Command  string
	ExitCode int
	Output   string
}

type CommandErrorMsg struct {
	Command string
	Error   error
	Output  string
}
