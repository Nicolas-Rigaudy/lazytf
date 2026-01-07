package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
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

	mainPanelErrorStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Red). // Red when error
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
	HasError  bool
	viewport  viewport.Model
	ready     bool
}

func NewMainPanel() MainPanelModel {
	vp := viewport.New(0, 0)
	vp.YPosition = 0
	vp.Style = lipgloss.NewStyle().UnsetMaxWidth() // Allow full width rendering
	return MainPanelModel{
		Content:   "Main Panel Content",
		Width:     0,
		Height:    0,
		IsFocused: false,
		viewport:  vp,
		ready:     false,
	}
}

func (m *MainPanelModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height

	// Account for border and padding (2 chars on each side for border, 1 for padding = 6 total width)
	// Account for title line if present (title + blank line + top/bottom border/padding = 7 total height)
	viewportWidth := width - 6
	viewportHeight := height - 7
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	if viewportWidth < 1 {
		viewportWidth = 1
	}

	m.viewport.Width = viewportWidth
	m.viewport.Height = viewportHeight
	m.ready = true
}

func (m *MainPanelModel) SetContent(content string) {
	m.Content = content
	// Set content directly without wrapping - let terminal handle line overflow
	// Wrapping is too expensive for real-time streaming output
	m.viewport.SetContent(content)
	// Auto-scroll to bottom when new content arrives
	m.viewport.GotoBottom()
}

func (m MainPanelModel) View() string {
	if !m.ready {
		return ""
	}

	// Build content with title
	var viewContent string
	if m.Title != "" {
		titleLine := mainPanelTitleStyle.Render(m.Title)
		viewContent = lipgloss.JoinVertical(lipgloss.Left, titleLine, "", m.viewport.View())
	} else {
		viewContent = m.viewport.View()
	}

	// Apply border style
	if m.HasError {
		return mainPanelErrorStyle.Width(m.Width).Height(m.Height).Render(viewContent)
	} else if m.IsFocused {
		return mainPanelFocusedStyle.Width(m.Width).Height(m.Height).Render(viewContent)
	} else {
		return mainPanelUnfocusedStyle.Width(m.Width).Height(m.Height).Render(viewContent)
	}
}
