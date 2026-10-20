package terraform

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type InitOptions struct {
	BackendConfigFile BackendVarFile
	Reconfigure       bool
	Upgrade           bool
	Input             bool
}

// buildInitArgs is the functional core: pure arg-building, no side effects.
func buildInitArgs(options InitOptions) []string {
	args := []string{"init"}

	if options.BackendConfigFile.Name != "" {
		args = append(args, fmt.Sprintf("-backend-config=%s", options.BackendConfigFile.FullPath))
	}

	if options.Reconfigure {
		args = append(args, "-reconfigure")
	}

	if options.Upgrade {
		args = append(args, "-upgrade")
	}

	if !options.Input {
		args = append(args, "-input=false")
	}

	return args
}

// RunInit is the imperative shell: builds args, then executes.
func RunInit(tool Tool, projectPath string, options InitOptions) tea.Cmd {
	args := buildInitArgs(options)
	result := tool.stream(args, projectPath, false)
	return result.Cmd
}
