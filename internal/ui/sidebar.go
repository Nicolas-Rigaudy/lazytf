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
	Title          string
	Items          []string
	SelectedIndex  int
	Width          int
	Height         int
	IsFocused      bool
	IsCollapsed    bool
	InitializedEnv string
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

// renderWithStyle applies the appropriate border style based on focus state
func (m SidebarModel) renderWithStyle(content string) string {
	style := sidebarUnfocusedStyle
	if m.IsFocused {
		style = sidebarFocusedStyle
	}
	return style.Width(m.Width).Height(m.Height).Render(content)
}

// renderItem renders a single item with selection highlighting
func (m SidebarModel) renderItem(index int, display string) string {
	if index == m.SelectedIndex {
		return highlightedItemStyle.Render(display)
	}
	return normalItemStyle.Render(display)
}

// isInitialized checks if the given item is the initialized environment
func (m SidebarModel) isInitialized(item string) bool {
	return m.InitializedEnv != "" && item == m.InitializedEnv
}

func (m SidebarModel) View() string {
	var parts []string

	if m.IsCollapsed {
		// Mini title - first 5 chars
		if m.Title != "" {
			miniTitle := m.Title
			if len(miniTitle) > 5 {
				miniTitle = miniTitle[:5]
			}
			parts = append(parts, sidebarTitleStyle.Render(miniTitle))
		}

		// Collapsed items - first letter with arrow and checkmark
		for i, item := range m.Items {
			suffix := " " // Space to match checkmark width
			if m.isInitialized(item) {
				suffix = "✅"
			}
			display := string(item[0]) + suffix
			// Use simpler styles without padding for collapsed view
			if i == m.SelectedIndex {
				parts = append(parts, highlightedItemStyle.Render(display))
			} else {
				parts = append(parts, normalItemStyle.Render(display))
			}
		}
	} else {
		// Full title
		if m.Title != "" {
			parts = append(parts, sidebarTitleStyle.Render("Project: "+m.Title))
		}

		// Expanded items - full name with initialized indicator
		for i, item := range m.Items {
			display := item
			if m.isInitialized(item) {
				display += " ✅ Initialized"
			}
			parts = append(parts, m.renderItem(i, display))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.renderWithStyle(content)
}
