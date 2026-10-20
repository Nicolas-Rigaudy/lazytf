package executor

import (
	"io"
	"os/exec"
	"slices"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// drain follows the ListenNext chain the way the UI does and returns the
// output lines and the final message.
func drain(t *testing.T, cmd tea.Cmd) (stdout, stderr []string, final tea.Msg) {
	t.Helper()
	deadline := time.After(10 * time.Second)
	for {
		msgCh := make(chan tea.Msg, 1)
		go func() { msgCh <- cmd() }()
		select {
		case <-deadline:
			t.Fatal("command did not finish within 10s")
		case msg := <-msgCh:
			out, ok := msg.(CommandOutputMsg)
			if !ok {
				return stdout, stderr, msg
			}
			if out.IsErr {
				stderr = append(stderr, out.Line)
			} else {
				stdout = append(stdout, out.Line)
			}
			cmd = out.ListenNext
		}
	}
}

func requireSh(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available")
	}
}

func TestExecuteStreamingSuccess(t *testing.T) {
	requireSh(t)
	res := ExecuteStreaming("sh", []string{"-c", "echo one; echo two; echo warn >&2"}, t.TempDir(), false)
	if res.StdinWriter != nil {
		t.Error("StdinWriter set without withStdin")
	}
	stdout, stderr, final := drain(t, res.Cmd)
	if !slices.Equal(stdout, []string{"one", "two"}) || !slices.Equal(stderr, []string{"warn"}) {
		t.Errorf("stdout %v stderr %v; want [one two] and [warn]", stdout, stderr)
	}
	if done, ok := final.(CommandCompletedMsg); !ok || done.ExitCode != 0 {
		t.Errorf("final = %#v, want CommandCompletedMsg with exit 0", final)
	}
}

func TestExecuteStreamingExitCode(t *testing.T) {
	requireSh(t)
	res := ExecuteStreaming("sh", []string{"-c", "echo boom >&2; exit 3"}, "", false)
	_, stderr, final := drain(t, res.Cmd)
	failed, ok := final.(CommandErrorMsg)
	if !ok || failed.ExitCode != 3 {
		t.Fatalf("final = %#v, want CommandErrorMsg with exit 3", final)
	}
	// All output must arrive before the failure message.
	if !slices.Equal(stderr, []string{"boom"}) {
		t.Errorf("stderr = %v, want [boom]", stderr)
	}
}

func TestExecuteStreamingMissingBinary(t *testing.T) {
	res := ExecuteStreaming("lazytf-no-such-binary", nil, "", false)
	_, _, final := drain(t, res.Cmd)
	failed, ok := final.(CommandErrorMsg)
	if !ok || failed.Error == nil {
		t.Fatalf("final = %#v, want CommandErrorMsg with an error", final)
	}
	if failed.Command != "lazytf-no-such-binary " {
		t.Errorf("Command = %q", failed.Command)
	}
}

func TestExecuteStreamingRunsInWorkingDir(t *testing.T) {
	requireSh(t)
	dir := t.TempDir()
	res := ExecuteStreaming("sh", []string{"-c", "pwd -P"}, dir, false)
	stdout, _, _ := drain(t, res.Cmd)
	want, _ := exec.Command("sh", "-c", "cd '"+dir+"' && pwd -P").Output()
	if len(stdout) != 1 || stdout[0]+"\n" != string(want) {
		t.Errorf("pwd = %v, want %q", stdout, want)
	}
}

func TestExecuteStreamingStdin(t *testing.T) {
	requireSh(t)
	res := ExecuteStreaming("sh", []string{"-c", "read answer; echo got $answer"}, "", true)
	if res.StdinWriter == nil {
		t.Fatal("StdinWriter is nil with withStdin")
	}
	if _, err := io.WriteString(res.StdinWriter, "yes\n"); err != nil {
		t.Fatal(err)
	}
	stdout, _, final := drain(t, res.Cmd)
	if !slices.Equal(stdout, []string{"got yes"}) {
		t.Errorf("stdout = %v, want [got yes]", stdout)
	}
	if _, ok := final.(CommandCompletedMsg); !ok {
		t.Errorf("final = %#v, want CommandCompletedMsg", final)
	}
}

func TestExecuteStreamingEnvAddsVariables(t *testing.T) {
	requireSh(t)
	t.Setenv("LAZYTF_KEPT", "inherited")
	res := ExecuteStreamingEnv("sh", []string{"-c", `echo "$LAZYTF_PROFILE $LAZYTF_KEPT"`}, t.TempDir(), false, []string{"LAZYTF_PROFILE=dev"})
	stdout, _, _ := drain(t, res.Cmd)
	if !slices.Equal(stdout, []string{"dev inherited"}) {
		t.Errorf("stdout %v, want the added variable and the inherited one", stdout)
	}
}
