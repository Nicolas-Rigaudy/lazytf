package ui

import (
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ModalType identifies which modal is currently active
type ModalType int

const (
	ModalNone    ModalType = iota
	ModalConfirm           // Generic confirmation modal (yes/no)
	ModalSelect            // Generic selection modal (pick from list)
	ModalError             // Error display modal
)

// ModalState holds the current modal state (pure data)
type ModalState struct {
	Type ModalType

	// For ModalConfirm
	Title     string
	Message   string
	OnConfirm func() tea.Msg // Called when user presses 'y' or Enter
	OnCancel  func() tea.Msg // Called when user presses 'n' or Esc

	// For ModalSelect
	Items    []string
	Selected int                             // Currently selected item index
	OnSelect func(selectedIndex int) tea.Msg // Called when user presses Enter

	// For ModalError
	ErrorText string
}

// ═══════════════════════════════════════════════════════════════════════════
// MODAL COMPONENT - Self-contained modal UI component
// ═══════════════════════════════════════════════════════════════════════════

// Modal is a self-contained modal component that manages its own state and rendering
// It follows the Bubble Tea component pattern (like SidebarModel, MainPanelModel)
type Modal struct {
	state ModalState
}

// NewModal creates a new modal component
func NewModal() Modal {
	return Modal{
		state: ModalState{
			Type: ModalNone,
		},
	}
}

// IsActive returns true if a modal is currently displayed
func (m Modal) IsActive() bool {
	return m.state.Type != ModalNone
}

// Show displays a modal with the given state
func (m *Modal) Show(state ModalState) {
	m.state = state
}

// Close closes the currently active modal
func (m *Modal) Close() {
	m.state = ModalState{Type: ModalNone}
}

// Update handles input for the modal
func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state.Type == ModalConfirm {
			switch msg.String() {
			case "y", "enter":
				if m.state.OnConfirm != nil {
					resultMsg := m.state.OnConfirm()
					m.state = ModalState{Type: ModalNone} // Close modal
					// Wrap the message in a Cmd
					return m, func() tea.Msg { return resultMsg }
				}
			case "n", "esc":
				m.state = ModalState{Type: ModalNone} // Close modal
				if m.state.OnCancel != nil {
					resultMsg := m.state.OnCancel()
					return m, func() tea.Msg { return resultMsg }
				}
				return m, nil
			}
		}
		if m.state.Type == ModalSelect {
			switch msg.String() {
			case "up", "k":
				if m.state.Selected > 0 {
					m.state.Selected--
				}
			case "down", "j":
				if m.state.Selected < len(m.state.Items)-1 {
					m.state.Selected++
				}
			case "enter":
				if m.state.OnSelect != nil {
					resultMsg := m.state.OnSelect(m.state.Selected)
					m.state = ModalState{Type: ModalNone} // Close modal
					// Wrap the message in a Cmd
					return m, func() tea.Msg { return resultMsg }
				}
			case "esc":
				m.state = ModalState{Type: ModalNone} // Close modal
				return m, nil
			}
		}
		if m.state.Type == ModalError {
			switch msg.String() {
			case "enter", "esc":
				m.state = ModalState{Type: ModalNone} // Close modal
				return m, nil
			}
		}
	}
	return m, nil
}

// View renders the modal based on its current state
func (m Modal) View(termWidth, termHeight int) string {
	switch m.state.Type {
	case ModalConfirm:
		return RenderConfirmModal(m.state, termWidth, termHeight)
	case ModalSelect:
		return RenderSelectModal(m.state, termWidth, termHeight)
	case ModalError:
		return RenderErrorModal(m.state, termWidth, termHeight)
	default:
		return ""
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// MODAL BUILDER - Reusable modal rendering component
// ═══════════════════════════════════════════════════════════════════════════

// ModalButton represents a button in a modal
type ModalButton struct {
	Label string         // Display text (e.g., "[y] Confirm")
	Color lipgloss.Color // Button color from theme
	Key   string         // Which key triggers this button (e.g., "y", "n")
}

// ModalBuilder is a flexible component for building any modal
// It composes different parts (title, content, buttons) into a centered modal
type ModalBuilder struct {
	Title       string         // Modal title (bold, colored)
	Content     string         // Main content (can be multi-line)
	Buttons     []ModalButton  // Action buttons at bottom
	Width       int            // Modal width in characters
	Height      int            // Modal height in lines
	BorderColor lipgloss.Color // Border color from theme
}

var modalTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(theme.Current.Mauve)
var modalContentStyle = lipgloss.NewStyle().Foreground(theme.Current.Text)
var modalButtonStyle = lipgloss.NewStyle().Padding(0, 1).Bold(true)

// Render builds and returns the final modal string
func (mb ModalBuilder) Render(termWidth, termHeight int) string {
	title := modalTitleStyle.Render(mb.Title)
	separator := strings.Repeat("─", mb.Width-4)
	content := strings.Join(strings.Split(mb.Content, "\n"), "\n")
	var buttonStrs []string
	for _, btn := range mb.Buttons {
		btnStyle := modalButtonStyle.Foreground(btn.Color)
		buttonStrs = append(buttonStrs, btnStyle.Render(btn.Label))
	}

	buttons := strings.Join(buttonStrs, "  ")

	modal := lipgloss.JoinVertical(lipgloss.Center, title, separator, "", content, "", buttons)

	modalBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(mb.BorderColor).
		Padding(1, 2).
		Width(mb.Width).
		Height(mb.Height).
		Render(modal)

	return lipgloss.Place(
		termWidth, termHeight,
		lipgloss.Center, lipgloss.Center,
		modalBox,
	)
}
