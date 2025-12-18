package executor

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// ExecuteStreaming runs a command and streams output line-by-line
// commandName: the executable to run (e.g., "terraform", "aws")
// args: command arguments
// workingDir: directory to run the command in (empty string for current dir)
func ExecuteStreaming(commandName string, args []string, workingDir string) tea.Cmd {
	// Create the command
	cmd := exec.Command(commandName, args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	cmdString := commandName + " " + strings.Join(args, " ")

	// Get stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return func() tea.Msg {
			return CommandErrorMsg{
				Command: cmdString,
				Error:   err,
				Output:  "Failed to create stdout pipe",
			}
		}
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return func() tea.Msg {
			return CommandErrorMsg{
				Command: cmdString,
				Error:   err,
				Output:  "Failed to create stderr pipe",
			}
		}
	}

	// Start the command (non-blocking)
	err = cmd.Start()
	if err != nil {
		return func() tea.Msg {
			return CommandErrorMsg{
				Command: cmdString,
				Error:   err,
				Output:  "Failed to start command",
			}
		}
	}

	// Create a channel to send messages back to Bubble Tea
	outputChannel := make(chan CommandOutputMsg)

	// WaitGroup to track when both stdout and stderr are done
	var wg sync.WaitGroup
	wg.Add(2) // We have 2 goroutines to wait for

	// Launch goroutines to read stdout and stderr
	go func() {
		defer wg.Done() // Signal when this goroutine finishes
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			outputChannel <- CommandOutputMsg{
				Line:  line,
				IsErr: false,
			}
		}
	}()

	go func() {
		defer wg.Done() // Signal when this goroutine finishes
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			outputChannel <- CommandOutputMsg{
				Line:  line,
				IsErr: true,
			}
		}
	}()

	// Launch a goroutine to wait for command completion
	go func() {
		cmd.Wait()           // Wait for command to finish
		wg.Wait()            // Wait for stdout and stderr goroutines to finish reading
		close(outputChannel) // NOW it's safe to close the channel
	}()

	// Return a listener that will read from the channel
	return listenToChannel(outputChannel, cmdString)
}

// listenToChannel creates a tea.Cmd that reads one message from a channel
// This function will be called repeatedly by the UI's Update function
func listenToChannel(ch chan CommandOutputMsg, cmdString string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			// Channel closed, command completed
			return CommandCompletedMsg{
				Command:  cmdString,
				ExitCode: 0,
				Output:   "",
			}
		}
		// Include the command to listen for the next message
		msg.ListenNext = listenToChannel(ch, cmdString)
		return msg
	}
}
