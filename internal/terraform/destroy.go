package terraform

import (
	"fmt"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// buildDestroyArgs is the functional core: pure arg-building, no side effects.
func buildDestroyArgs(varFile VarFile) []string {
	args := []string{"destroy"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	return args
}

// RunDestroy executes `terraform destroy -var-file=...` (imperative shell).
// Always interactive - returns InteractiveResult with stdin writer for sending confirmation.
func RunDestroy(projectPath string, varFile VarFile) InteractiveResult {
	args := buildDestroyArgs(varFile)
	result := executor.ExecuteStreaming("terraform", args, projectPath, true)
	return InteractiveResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
