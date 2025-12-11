package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	sidebarTitleStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Mauve).
				Bold(true).
				Underline(true).
				PaddingBottom(1)

	sidebarFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Blue). // Blue when focused
				Padding(1)

	sidebarUnfocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Current.Surface2). // Gray when unfocused
				Padding(1)

	highlightedItemStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Crust).
				Background(theme.Current.Blue).
				Bold(true).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(theme.Current.Text).
			Padding(0, 1)
)

type SidebarModel struct {
	Title         string
	Items         []string
	SelectedIndex int
	Width         int
	Height        int
	IsFocused     bool
}

func NewSidebar(items ...string) SidebarModel {
	return SidebarModel{
		Items:         items,
		SelectedIndex: 0,
		Width:         0,
		Height:        0,
		IsFocused:     false,
	}
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
		case "down", "j":
			if m.SelectedIndex < len(m.Items)-1 {
				m.SelectedIndex++
			}
		}
	}
	return m, nil
}

func (m SidebarModel) View() string {
	var parts []string
	if m.Title != "" {
		parts = append(parts, sidebarTitleStyle.Render("🎯 Project: "+m.Title))
		parts = append(parts, "") // Add a blank line after title
	}

	items := []string{}
	for i, item := range m.Items {
		if i == m.SelectedIndex {
			items = append(items, highlightedItemStyle.Render(item))
		} else {
			items = append(items, normalItemStyle.Render(item))
		}
	}

	parts = append(parts, items...)

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	if m.IsFocused {
		return sidebarFocusedStyle.Width(m.Width).Height(m.Height).Render(content)
	}
	return sidebarUnfocusedStyle.Width(m.Width).Height(m.Height).Render(content)
}
