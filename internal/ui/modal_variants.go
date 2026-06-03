package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════════════════
// MODAL VARIANTS - Ready-to-use modal types
// ═══════════════════════════════════════════════════════════════════════════
//
// This file contains convenience functions for creating common modal types.
// Each function uses the ModalBuilder to construct a specific kind of modal.
//
// ═══════════════════════════════════════════════════════════════════════════

func RenderConfirmModal(state ModalState, termWidth, termHeight int) string {
	confirmLabel := state.ConfirmLabel
	if confirmLabel == "" {
		confirmLabel = "Yes"
	}
	cancelLabel := state.CancelLabel
	if cancelLabel == "" {
		cancelLabel = "No"
	}

	builder := ModalBuilder{
		Title:   state.Title,
		Content: state.Message,
		Buttons: []ModalButton{
			{Label: "[Enter/y] " + confirmLabel, Color: theme.Current.Green, Key: "y"},
			{Label: "[ESC/n] " + cancelLabel, Color: theme.Current.Red, Key: "n"},
		},
		Width:       60,
		Height:      14,
		BorderColor: theme.Current.Blue,
	}

	return builder.Render(termWidth, termHeight)
}

func RenderSelectModal(state ModalState, termWidth, termHeight int) string {
	content := ""
	for i, item := range state.Items {
		if i == state.Selected {
			content += "› " + item + "\n"

		} else {
			content += "  " + item + "\n"
		}
	}

	builder := ModalBuilder{
		Title:   state.Title,
		Content: content,
		Buttons: []ModalButton{
			{Label: "[Enter] Select", Color: theme.Current.Green, Key: "enter"},
			{Label: "[ESC] Cancel", Color: theme.Current.Red, Key: "esc"},
		},
		Width:       50,
		Height:      len(state.Items) + 8, // Dynamic height based on items
		BorderColor: theme.Current.Blue,
	}

	return builder.Render(termWidth, termHeight)
}

// RenderErrorModal renders an error display modal
func RenderErrorModal(state ModalState, termWidth, termHeight int) string {
	title := state.Title
	if title == "" {
		title = "❌ Error"
	}

	builder := ModalBuilder{
		Title:   title,
		Content: state.ErrorText,
		Buttons: []ModalButton{
			{Label: "[Enter] OK", Color: theme.Current.Red, Key: "enter"},
		},
		Width:       70,
		Height:      12,
		BorderColor: theme.Current.Red, // Red border for errors
	}

	return builder.Render(termWidth, termHeight)
}

func RenderHelpModal(state ModalState, termWidth, termHeight int) string {
	var rows [][]string

	for i, group := range state.KeyBindingGroups {
		if i > 0 {
			rows = append(rows, []string{"", ""}) // spacer between groups
		}
		rows = append(rows, []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Current.Mauve).Render(group.Title),
			"",
		})
		for _, kb := range group.KeyBindings {
			var keyStr, descStr string
			if kb.Enabled {
				keyStr = lipgloss.NewStyle().Foreground(theme.Current.Green).Render(kb.Key)
				descStr = kb.Description
			} else {
				keyStr = lipgloss.NewStyle().Foreground(theme.Current.Overlay0).Render(kb.Key)
				descStr = lipgloss.NewStyle().Foreground(theme.Current.Overlay0).Render(kb.Description)
			}
			rows = append(rows, []string{keyStr, descStr})
		}
	}

	totalRows := len(rows)
	builder := ModalBuilder{
		Title: "📖 Help - Key Bindings",
		Rows:  rows,
		Buttons: []ModalButton{
			{Label: "[ESC] Close", Color: theme.Current.Blue, Key: "esc"},
		},
		Width:       60,
		Height:      totalRows + 10,
		BorderColor: theme.Current.Blue,
	}

	return builder.Render(termWidth, termHeight)
}
