package terraform

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// InteractiveResult contains both the tea.Cmd and stdin writer for commands that need user confirmation
type InteractiveResult struct {
	Cmd         tea.Cmd
	StdinWriter io.WriteCloser
}

// RunApply executes `terraform apply -var-file=...`
// Returns both the command and a stdin writer to send "yes" confirmation later
func RunApply(projectPath string, varFile VarFile, planFile string) InteractiveResult {
	args := []string{"apply"}
	var withStdin bool

	if planFile != "" {
		args = append(args, planFile)
		withStdin = false // No need for stdin confirmation when using a plan file
	} else {
		if varFile.FullPath != "" {
			args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
		}
		withStdin = true // Will need stdin confirmation for interactive apply
	}

	result := executor.ExecuteStreaming("terraform", args, projectPath, withStdin)
	return InteractiveResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
