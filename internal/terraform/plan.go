package terraform

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// buildPlanArgs is the functional core: pure arg-building, no side effects.
// Kept unexported (implementation detail) but reachable from same-package tests.
func buildPlanArgs(varFile VarFile, planFile string) []string {
	// -input=false: a missing variable fails at once instead of waiting on
	// a stdin that is not wired.
	args := []string{"plan", "-input=false"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	args = append(args, "-out="+planFile)

	return args
}

// RunPlan is the imperative shell: builds args, then executes.
func RunPlan(bin, projectPath string, varFile VarFile) tea.Cmd {
	planFile := PlanFilePath(projectPath, varFile.EnvName)
	// Plans can embed sensitive values: keep the directory private.
	// A mkdir failure surfaces as Terraform's own -out error.
	_ = os.MkdirAll(filepath.Dir(planFile), 0o700)
	args := buildPlanArgs(varFile, planFile)
	result := executor.ExecuteStreaming(Binary(bin), args, projectPath, false)
	return result.Cmd
}

// PlanFilePath returns where the saved plan for a project and env lives.
// It is keyed by absolute project path and env name, so two envs of one
// project, or two projects sharing a directory name, never share a plan.
func PlanFilePath(projectPath, envName string) string {
	return planFilePath(planCacheDir(), projectPath, envName)
}

func planCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "lazytf", "plans")
}

func planFilePath(cacheDir, projectPath, envName string) string {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		abs = filepath.Clean(projectPath)
	}
	sum := sha256.Sum256([]byte(abs))
	if envName == "" {
		envName = "default"
	}
	name := fmt.Sprintf("%s-%x-%s.tfplan", filepath.Base(abs), sum[:6], filepath.Base(envName))
	return filepath.Join(cacheDir, name)
}

// LastInitTime returns when `terraform init` last wrote the backend config
// (.terraform/terraform.tfstate). ok is false when it cannot be determined.
func LastInitTime(projectPath string) (t time.Time, ok bool) {
	info, err := os.Stat(filepath.Join(projectPath, ".terraform", "terraform.tfstate"))
	if err != nil {
		return time.Time{}, false
	}
	return info.ModTime(), true
}

// PlanOlderThanInit reports whether a saved plan predates the last init,
// in which case it may target a different backend than the current one.
func PlanOlderThanInit(projectPath string, planModTime time.Time) bool {
	initTime, ok := LastInitTime(projectPath)
	return ok && planModTime.Before(initTime)
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
