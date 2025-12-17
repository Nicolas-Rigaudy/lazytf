package ui

import (
	"time"

	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerContainerStyle = lipgloss.NewStyle().
				Padding(0, 1)

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
		Height: 3,
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

	h.Height = lipgloss.Height(content) + 1 // +1 for padding
	return headerContainerStyle.Width(h.Width).Height(h.Height).Render(content)
}
