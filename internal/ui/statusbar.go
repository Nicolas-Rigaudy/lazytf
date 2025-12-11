package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
		Foreground(theme.Current.Text).
		Background(theme.Current.Surface0).
		Padding(0, 2)
)

type StatusBarModel struct {
	Text  string
	Width int
}

func NewStatusBar() StatusBarModel {
	return StatusBarModel{
		Text:  "q: quit | tab: switch panels | ↑↓/jk: navigate",
		Width: 0,
	}
}

func (s StatusBarModel) View() string {
	return statusBarStyle.Width(s.Width).Render(s.Text)
}

func (s *StatusBarModel) SetText(text string) {
	s.Text = text
}
