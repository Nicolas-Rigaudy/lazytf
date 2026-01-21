package terraform

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// ApplyResult contains both the tea.Cmd and stdin writer for sending confirmation
type ApplyResult struct {
	Cmd         tea.Cmd
	StdinWriter io.WriteCloser
}

// RunApply executes `terraform apply -var-file=...`
// Returns both the command and a stdin writer to send "yes" confirmation later
func RunApply(projectPath string, varFile VarFile) ApplyResult {
	args := []string{"apply"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	result := executor.ExecuteStreaming("terraform", args, projectPath, true)
	return ApplyResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
