package theme

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleEnv selects the colour style: "catppuccin" (default), "ansi-a" or "ansi-b".
// The ANSI styles use only the terminal's 16 theme colours, so they follow the
// terminal theme, and paint no backgrounds.
const StyleEnv = "LAZYTF_STYLE"

// Roles names what each colour is used for, so a style can map them freely.
type Roles struct {
	Name           string
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

	Selected       lipgloss.Style // the selected plan row
	SelectedItem   lipgloss.Style // the selected sidebar row
	SelectedMarker string         // drawn before the selected sidebar row, blanks before the others
	TitleBar       lipgloss.Style
	StatusBar      lipgloss.Style
	ConfirmPrompt  lipgloss.Style
}

// Active is the style in use, chosen from StyleEnv at start-up.
var Active = ForName(os.Getenv(StyleEnv))

// ForName returns the named style; unknown names give the default.
func ForName(name string) Roles {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "ansi-a":
		return ansi("ansi-a", "2", "6", "6", lipgloss.NewStyle().Reverse(true).Bold(true), "")
	case "ansi-b":
		return ansi("ansi-b", "4", "5", "5", lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true), "▶ ")
	}
	return catppuccin()
}

func catppuccin() Roles {
	c := Current
	return Roles{
		Name: "catppuccin", ActiveBorder: c.Blue, InactiveBorder: c.Surface2,
		Accent: c.Mauve, Hint: c.Green, Text: c.Text, Subtle: c.Subtext0, Muted: c.Overlay0,
		Add: c.Green, Change: c.Yellow, Destroy: c.Red,
		Selected:      lipgloss.NewStyle().Background(c.Surface1),
		SelectedItem:  lipgloss.NewStyle().Foreground(c.Crust).Background(c.Blue).Bold(true),
		TitleBar:      lipgloss.NewStyle().Background(c.Mauve).Foreground(c.Crust).Bold(true),
		StatusBar:     lipgloss.NewStyle().Foreground(c.Text).Background(c.Surface0),
		ConfirmPrompt: lipgloss.NewStyle().Foreground(c.Base).Background(c.Mauve).Bold(true),
	}
}

// ansi builds an ANSI-slot style; the two share everything but the border,
// accent and selected-row looks.
func ansi(name, border, accent, hint string, selected lipgloss.Style, marker string) Roles {
	fg := func(c string) lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Bold(true) }
	return Roles{
		Name: name, ActiveBorder: lipgloss.Color(border), InactiveBorder: "8",
		Accent: lipgloss.Color(accent), Hint: lipgloss.Color(hint), Text: "", Subtle: "8", Muted: "8",
		Add: "2", Change: "3", Destroy: "1",
		Selected: selected, SelectedItem: selected, SelectedMarker: marker,
		TitleBar:      fg(accent),
		StatusBar:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		ConfirmPrompt: fg(accent),
	}
}
