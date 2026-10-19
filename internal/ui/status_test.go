package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/executor"
)

func TestCommandFailureText(t *testing.T) {
	tests := []struct {
		name string
		msg  executor.CommandErrorMsg
		want string
	}{
		{
			name: "non-zero exit",
			msg:  executor.CommandErrorMsg{Command: "terraform plan -input=false -var-file=/a/dev.tfvars", ExitCode: 1, Error: errors.New("exit status 1")},
			want: "✗ terraform plan failed: exit code 1",
		},
		{
			name: "could not start",
			msg:  executor.CommandErrorMsg{Command: "terraform init", Error: errors.New(`exec: "terraform": executable file not found in $PATH`), Output: "Failed to start command"},
			want: `✗ terraform init failed: exec: "terraform": executable file not found in $PATH`,
		},
		{
			name: "killed by signal",
			msg:  executor.CommandErrorMsg{Command: "terraform apply", ExitCode: -1, Error: errors.New("signal: killed")},
			want: "✗ terraform apply failed: signal: killed",
		},
		{
			name: "no error value",
			msg:  executor.CommandErrorMsg{Command: "terraform", Output: "Failed to create stdout pipe"},
			want: "✗ terraform failed: Failed to create stdout pipe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandFailureText(tt.msg); got != tt.want {
				t.Errorf("commandFailureText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatusTextShowsAndClearsFailure(t *testing.T) {
	var m Model
	m.lastFailure = "✗ terraform plan failed: exit code 1"
	if got := m.buildStatusText(); !strings.HasPrefix(got, m.lastFailure) {
		t.Errorf("buildStatusText() = %q, want it to start with the failure", got)
	}

	m.prepareCommandExecution()
	if m.lastFailure != "" || strings.Contains(m.statusBar.Text, "failed") {
		t.Errorf("after starting a new command, status = %q, want the failure cleared", m.statusBar.Text)
	}
}

func TestStatusTextShowsAndClearsNotice(t *testing.T) {
	var m Model
	m.SetNotice("")
	if strings.Contains(m.statusBar.Text, "Removed") {
		t.Errorf("empty notice shown: %q", m.statusBar.Text)
	}

	note := "Removed 2 old plan files from /tmp (v0.1 location)"
	m.SetNotice(note)
	if !strings.HasPrefix(m.statusBar.Text, note) {
		t.Errorf("status = %q, want it to start with the notice", m.statusBar.Text)
	}

	m.prepareCommandExecution()
	if m.notice != "" || strings.Contains(m.statusBar.Text, note) {
		t.Errorf("after starting a new command, status = %q, want the notice cleared", m.statusBar.Text)
	}
}
