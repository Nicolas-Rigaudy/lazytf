package ui

import (
	"time"

	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Surface2).
				Padding(1, 2)

	headerLabelStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Mauve).
				Bold(true)

	headerValueStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Text)

	headerSuccessStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Green)

	headerErrorStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Red)
)

type HeaderModel struct {
	Width  int
	Height int
}

type InfoHeaderData struct {
	ProjectName     string
	EnvName         string
	IsInitialized   bool
	LastCommand     string
	LastCommandTime time.Time
}

func NewHeader() HeaderModel {
	return HeaderModel{
		Width:  0,
		Height: 5, // Border (2) + Padding (2) + Content (1 line) = 5
	}
}

func (h HeaderModel) View(data InfoHeaderData) string {
	line1 := lipgloss.JoinHorizontal(lipgloss.Left,
		headerLabelStyle.Render("Project: "),
		headerValueStyle.Render(data.ProjectName),
		"  ",
		headerLabelStyle.Render("Env: "),
		headerValueStyle.Render(data.EnvName),
		"  ",
		func() string {
			if data.IsInitialized {
				return headerSuccessStyle.Render("✅ Initialized")
			}
			return headerErrorStyle.Render("❌ Not Initialized")
		}(),
	)
	line2 := ""
	if data.LastCommand != "" {
		line2 = lipgloss.NewStyle().Foreground(theme.Current.Subtext0).Render(
			"Last: " + data.LastCommand + " (" + data.LastCommandTime.Format("15:04:05") + ")",
		)
	}

	// Join lines vertically
	content := line1
	if line2 != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, line1, line2)
	}

	// Render with border and padding
	rendered := headerContainerStyle.Width(h.Width).Render(content)
	h.Height = lipgloss.Height(rendered)
	return rendered
}
