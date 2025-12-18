package terraform

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

type InitOptions struct {
	BackendConfigFile BackendVarFile
	Reconfigure       bool
	Upgrade           bool
	Input             bool
}

func RunInit(projectPath string, options InitOptions) tea.Cmd {
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
	return executor.ExecuteStreaming("terraform", args, projectPath)
}
