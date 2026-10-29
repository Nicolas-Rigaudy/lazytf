package theme

import "github.com/charmbracelet/lipgloss"

// Roles names what each colour is used for. Colours are the terminal's 16 ANSI
// slots and no style paints a background, so LazyTF follows the terminal theme.
type Roles struct {
	ActiveBorder   lipgloss.Color // border of the focused pane, modal borders
	InactiveBorder lipgloss.Color
	Accent         lipgloss.Color // titles and labels
	Hint           lipgloss.Color // key hints
	Text           lipgloss.Color
	Subtle         lipgloss.Color // secondary text
	Muted          lipgloss.Color // disabled text
	Add            lipgloss.Color
	Change         lipgloss.Color
	Destroy        lipgloss.Color

	Selected      lipgloss.Style // the selected plan row
	SelectedItem  lipgloss.Style // the selected sidebar row
	TitleBar      lipgloss.Style
	StatusBar     lipgloss.Style
	ConfirmPrompt lipgloss.Style
}

// Active is the lazygit-like style: green active border, reverse-video
// selection, cyan accents, bright black for muted text.
var Active = Roles{
	ActiveBorder: "2", InactiveBorder: "8",
	Accent: "6", Hint: "6", Text: "", Subtle: "8", Muted: "8",
	Add: "2", Change: "3", Destroy: "1",
	Selected:      lipgloss.NewStyle().Reverse(true).Bold(true),
	SelectedItem:  lipgloss.NewStyle().Reverse(true).Bold(true),
	TitleBar:      lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true),
	StatusBar:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	ConfirmPrompt: lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true),
}
