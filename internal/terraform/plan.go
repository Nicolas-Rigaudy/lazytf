package terraform

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

func RunPlan(projectPath string, varFile VarFile) tea.Cmd {
	args := []string{"plan"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	args = append(args, "-out="+PlanFilePath(projectPath))

	result := executor.ExecuteStreaming("terraform", args, projectPath, false)
	return result.Cmd
}

func PlanFilePath(projectPath string) string {
	return fmt.Sprintf("/tmp/lazytf-%s.tfplan", filepath.Base(projectPath))
}
