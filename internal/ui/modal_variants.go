package ui

import (
	"github.com/Nicolas-Rigaudy/lazytf/internal/ui/theme"
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

	builder := ModalBuilder{
		Title:   state.Title,
		Content: state.Message,
		Buttons: []ModalButton{
			{Label: "[y] Yes", Color: theme.Current.Green, Key: "y"},
			{Label: "[n] No", Color: theme.Current.Red, Key: "n"},
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
