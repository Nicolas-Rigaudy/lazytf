package terraform

import (
	"fmt"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// buildPlanArgs is the functional core: pure arg-building, no side effects.
// Kept unexported (implementation detail) but reachable from same-package tests.
func buildPlanArgs(projectPath string, varFile VarFile) []string {
	args := []string{"plan"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	args = append(args, "-out="+PlanFilePath(projectPath))

	return args
}

// RunPlan is the imperative shell: builds args, then executes.
func RunPlan(projectPath string, varFile VarFile) tea.Cmd {
	args := buildPlanArgs(projectPath, varFile)
	result := executor.ExecuteStreaming("terraform", args, projectPath, false)
	return result.Cmd
}

func PlanFilePath(projectPath string) string {
	return fmt.Sprintf("/tmp/lazytf-%s.tfplan", filepath.Base(projectPath))
}

func PlanFileAge(modTime time.Time) string {
	duration := time.Since(modTime)
	if duration < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}
