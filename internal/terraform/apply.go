package terraform

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
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
	// -input=false turns off variable prompts only; Terraform still asks for
	// approval, so the "yes" sent on stdin can never fill a missing variable.
	args := []string{"apply", "-input=false"}
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
func RunApply(tool Tool, projectPath string, varFile VarFile, planFile string) InteractiveResult {
	args, withStdin := buildApplyArgs(varFile, planFile)
	result := tool.stream(args, projectPath, withStdin)
	return InteractiveResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
