package executor

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// ExecuteResult contains the tea.Cmd and an optional stdin writer
type ExecuteResult struct {
	Cmd         tea.Cmd
	StdinWriter io.WriteCloser // nil if command doesn't need stdin
}

// ExecuteStreaming runs a command and streams output line-by-line
// commandName: the executable to run (e.g., "terraform", "aws")
// args: command arguments
// workingDir: directory to run the command in (empty string for current dir)
// withStdin: if true, returns a stdin writer for sending input to the command
func ExecuteStreaming(commandName string, args []string, workingDir string, withStdin bool) ExecuteResult {
	// Create the command
	cmd := exec.Command(commandName, args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	cmdString := commandName + " " + strings.Join(args, " ")

	// Get stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return ExecuteResult{
			Cmd: func() tea.Msg {
				return CommandErrorMsg{
					Command: cmdString,
					Error:   err,
					Output:  "Failed to create stdout pipe",
				}
			},
			StdinWriter: nil,
		}
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return ExecuteResult{
			Cmd: func() tea.Msg {
				return CommandErrorMsg{
					Command: cmdString,
					Error:   err,
					Output:  "Failed to create stderr pipe",
				}
			},
			StdinWriter: nil,
		}
	}

	// Get stdin pipe if requested - must be done BEFORE cmd.Start()
	var stdinPipe io.WriteCloser
	if withStdin {
		stdinPipe, err = cmd.StdinPipe()
		if err != nil {
			return ExecuteResult{
				Cmd: func() tea.Msg {
					return CommandErrorMsg{
						Command: cmdString,
						Error:   err,
						Output:  "Failed to create stdin pipe",
					}
				},
				StdinWriter: nil,
			}
		}
	}

	// Start the command (non-blocking)
	err = cmd.Start()
	if err != nil {
		return ExecuteResult{
			Cmd: func() tea.Msg {
				return CommandErrorMsg{
					Command: cmdString,
					Error:   err,
					Output:  "Failed to start command",
				}
			},
			StdinWriter: nil,
		}
	}

	// Create a channel to send messages back to Bubble Tea
	outputChannel := make(chan tea.Msg)

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
		// CRITICAL: Wait for stdout/stderr to finish reading FIRST
		wg.Wait() // Wait for stdout and stderr goroutines to finish reading

		// NOW wait for the command to finish
		err := cmd.Wait()

		// Send completion/error message AFTER all output is captured
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				outputChannel <- CommandErrorMsg{
					Command:  cmdString,
					Error:    err,
					Output:   "Command failed",
					ExitCode: exitCode,
				}
			} else {
				outputChannel <- CommandErrorMsg{
					Command:  cmdString,
					Error:    err,
					Output:   "Command failed",
					ExitCode: -1,
				}
			}
		} else {
			outputChannel <- CommandCompletedMsg{
				Command:  cmdString,
				ExitCode: 0,
				Output:   "Command completed successfully",
			}
		}

		close(outputChannel) // NOW it's safe to close the channel
	}()

	// Return a listener that will read from the channel and the stdin writer
	return ExecuteResult{
		Cmd:         listenToChannel(outputChannel, cmdString),
		StdinWriter: stdinPipe, // Will be nil if withStdin was false
	}
}

// listenToChannel creates a tea.Cmd that reads one message from a channel
// This function will be called repeatedly by the UI's Update function
func listenToChannel(ch chan tea.Msg, cmdString string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			// Channel closed - this won't be used anymore since we'll send completion msgs directly
			return CommandCompletedMsg{
				Command:  cmdString,
				ExitCode: 0,
				Output:   "",
			}
		}
		if outputMsg, ok := msg.(CommandOutputMsg); ok {
			outputMsg.ListenNext = listenToChannel(ch, cmdString)
			return outputMsg
		}
		return msg
	}
}
