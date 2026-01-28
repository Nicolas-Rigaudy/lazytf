package ui

import (
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

	confirmPromptStyle = lipgloss.NewStyle().
				Foreground(theme.Current.Base).
				Background(theme.Current.Mauve).
				Bold(true).
				Padding(0, 1)
)

const (
	ViewModeRaw     = "raw"
	ViewModeSummary = "planSummary"
)

type MainPanelModel struct {
	Title         string
	Content       string
	Width         int
	Height        int
	IsFocused     bool
	HasError      bool
	viewport      viewport.Model
	ready         bool
	ConfirmPrompt string
	ViewMode      string
	PlanSummary   *PlanSummaryViewModel
}

func NewMainPanel() MainPanelModel {
	vp := viewport.New(0, 0)
	vp.KeyMap = viewport.DefaultKeyMap()
	vp.YPosition = 0
	vp.Style = lipgloss.NewStyle().UnsetMaxWidth() // Allow full width rendering
	vp.SetHorizontalStep(6)
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

	if m.PlanSummary != nil {
		m.PlanSummary.SetSize(viewportWidth, viewportHeight)
	}

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
	if m.ViewMode == ViewModeSummary && m.PlanSummary != nil {
		if m.Title != "" {
			titleLine := mainPanelTitleStyle.Render(m.Title)
			viewContent = lipgloss.JoinVertical(lipgloss.Left, titleLine, "", m.PlanSummary.View())
		} else {
			viewContent = m.PlanSummary.View()
		}
	} else if m.Title != "" {
		titleLine := mainPanelTitleStyle.Render(m.Title)
		viewContent = lipgloss.JoinVertical(lipgloss.Left, titleLine, "", m.viewport.View())
	} else {
		viewContent = m.viewport.View()
	}

	// Add confirm prompt if set
	if m.ConfirmPrompt != "" {
		separator := strings.Repeat("─", m.viewport.Width)
		prompt := confirmPromptStyle.Render(m.ConfirmPrompt)
		viewContent = lipgloss.JoinVertical(lipgloss.Left, viewContent, separator, prompt)
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

func (m *MainPanelModel) Update(msg tea.Msg) (MainPanelModel, tea.Cmd) {
	var cmd tea.Cmd

	// Handle view mode toggle
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "v":
			if m.ViewMode == ViewModeSummary {
				m.ViewMode = ViewModeRaw
			} else if m.PlanSummary != nil {
				m.ViewMode = ViewModeSummary
			}
			return *m, nil
		}
	}

	//Delegate to appropriate view updater
	if m.ViewMode == ViewModeSummary && m.PlanSummary != nil {
		m.PlanSummary.Update(msg)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return *m, cmd
}
