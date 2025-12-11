package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

var (
	mainPanelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Blue). // Blue when focused
				Padding(1)

	mainPanelUnfocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Surface2). // Gray when unfocused
				Padding(1)

	mainPanelTitleStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Mauve).
				Bold(true).
				PaddingBottom(1)
)

type MainPanelModel struct {
	Title     string
	Content   string
	Width     int
	Height    int
	IsFocused bool
}

func NewMainPanel() MainPanelModel {
	return MainPanelModel{
		Content:   "Main Panel Content",
		Width:     0,
		Height:    0,
		IsFocused: false,
	}
}

func (m MainPanelModel) View() string {
	// Build content with optional title at the top (inside border)
	var content string
	if m.Title != "" {
		titleLine := mainPanelTitleStyle.Render(m.Title)
		content = lipgloss.JoinVertical(lipgloss.Left, titleLine, "", m.Content)
	} else {
		content = m.Content
	}

	// Apply border style based on focus
	if m.IsFocused {
		return mainPanelFocusedStyle.Width(m.Width).Height(m.Height).Render(content)
	}
	return mainPanelUnfocusedStyle.Width(m.Width).Height(m.Height).Render(content)
}
