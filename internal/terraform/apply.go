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

// buildApplyArgs is the functional core: pure arg-building, no side effects.
// Returns the args plus withStdin, since the plan-file branch also decides
// whether interactive confirmation (stdin) is needed.
func buildApplyArgs(varFile VarFile, planFile string) ([]string, bool) {
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

	return args, withStdin
}

// RunApply executes `terraform apply -var-file=...` (imperative shell).
// Returns both the command and a stdin writer to send "yes" confirmation later
func RunApply(projectPath string, varFile VarFile, planFile string) InteractiveResult {
	args, withStdin := buildApplyArgs(varFile, planFile)
	result := executor.ExecuteStreaming("terraform", args, projectPath, withStdin)
	return InteractiveResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
