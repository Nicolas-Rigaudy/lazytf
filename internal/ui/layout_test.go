package ui

import (
	"strings"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
	"github.com/Nicolas-Rigaudy/lazytf/internal/terraform"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Bubble Tea cuts every line at the terminal width, so a line that is too wide
// loses its right border.
func TestViewFitsTerminalWidth(t *testing.T) {
	for _, width := range []int{80, 110, 160} {
		m := NewModel([]terraform.Project{{Name: "app", Path: t.TempDir()}}, terraform.ModeMultiProject, config.Layout{})
		updated, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 30})
		for i, line := range strings.Split(updated.View(), "\n") {
			if w := lipgloss.Width(line); w > width {
				t.Errorf("width %d: line %d is %d wide: %q", width, i, w, line)
			}
		}
	}
}
