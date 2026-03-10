package terraform

import (
	"fmt"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

// RunDestroy executes `terraform destroy -var-file=...`
// Always interactive - returns InteractiveResult with stdin writer for sending confirmation
func RunDestroy(projectPath string, varFile VarFile) InteractiveResult {
	args := []string{"destroy"}

	if varFile.FullPath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile.FullPath))
	}

	result := executor.ExecuteStreaming("terraform", args, projectPath, true)
	return InteractiveResult{
		Cmd:         result.Cmd,
		StdinWriter: result.StdinWriter,
	}
}
