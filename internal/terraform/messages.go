package terraform

type CommandOutputMsg struct {
	Line  string
	IsErr bool // true for stderr
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
