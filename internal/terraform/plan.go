package terraform

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

func RunPlan(projectPath string, varFile VarFile) tea.Cmd {
	args := []string{"plan"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	return executor.ExecuteStreaming("terraform", args, projectPath)
}
