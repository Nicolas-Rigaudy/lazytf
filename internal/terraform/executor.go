package terraform

import (
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func executeCommand(projectPath string, args []string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("terraform", args...)
		cmd.Dir = projectPath
		outputBytes, err := cmd.CombinedOutput()
		output := string(outputBytes)
		if err != nil {
			return CommandErrorMsg{
				Command: "terraform " + strings.Join(args, " "),
				Error:   err,
				Output:  output,
			}
		}
		return CommandCompletedMsg{
			Command:  "terraform " + strings.Join(args, " "),
			ExitCode: 0,
			Output:   output,
		}
	}
}
