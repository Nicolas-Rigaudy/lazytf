package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleBarStyle = lipgloss.NewStyle().
		Background(theme.Current.Mauve).
		Foreground(theme.Current.Crust).
		Bold(true).
		Padding(0, 2).
		Align(lipgloss.Center)
)

type TitleBarModel struct {
	AppName string
	Context string
	Width   int
}

func NewTitleBar() TitleBarModel {
	return TitleBarModel{
		AppName: "LazyTF",
		Context: "Terraform Commander",
		Width:   0,
	}
}

func (t TitleBarModel) View() string {
	title := t.AppName
	if t.Context != "" {
		title = title + " - " + t.Context
	}
	return titleBarStyle.Width(t.Width).Render(title)
}
