package theme

import "github.com/charmbracelet/lipgloss"

// Catppuccin Macchiato color palette
// Reference: https://github.com/catppuccin/catppuccin
type CatppuccinMacchiato struct {
	// Base colors
	Base       lipgloss.Color // #24273a - Background
	Mantle     lipgloss.Color // #1e2030 - Darker background
	Crust      lipgloss.Color // #181926 - Darkest background
	Surface0   lipgloss.Color // #363a4f - Surface
	Surface1   lipgloss.Color // #494d64 - Surface
	Surface2   lipgloss.Color // #5b6078 - Surface
	Overlay0   lipgloss.Color // #6e738d - Overlay
	Overlay1   lipgloss.Color // #8087a2 - Overlay
	Overlay2   lipgloss.Color // #939ab7 - Overlay
	Subtext0   lipgloss.Color // #a5adcb - Muted text
	Subtext1   lipgloss.Color // #b8c0e0 - Muted text
	Text       lipgloss.Color // #cad3f5 - Primary text

	// Accent colors
	Rosewater  lipgloss.Color // #f4dbd6
	Flamingo   lipgloss.Color // #f0c6c6
	Pink       lipgloss.Color // #f5bde6
	Mauve      lipgloss.Color // #c6a0f6 - Purple
	Red        lipgloss.Color // #ed8796
	Maroon     lipgloss.Color // #ee99a0
	Peach      lipgloss.Color // #f5a97f
	Yellow     lipgloss.Color // #eed49f
	Green      lipgloss.Color // #a6da95
	Teal       lipgloss.Color // #8bd5ca
	Sky        lipgloss.Color // #91d7e3
	Sapphire   lipgloss.Color // #7dc4e4
	Blue       lipgloss.Color // #8aadf4
	Lavender   lipgloss.Color // #b7bdf8
}

// NewCatppuccinMacchiato creates a new Catppuccin Macchiato theme
func NewCatppuccinMacchiato() *CatppuccinMacchiato {
	return &CatppuccinMacchiato{
		Base:      lipgloss.Color("#24273a"),
		Mantle:    lipgloss.Color("#1e2030"),
		Crust:     lipgloss.Color("#181926"),
		Surface0:  lipgloss.Color("#363a4f"),
		Surface1:  lipgloss.Color("#494d64"),
		Surface2:  lipgloss.Color("#5b6078"),
		Overlay0:  lipgloss.Color("#6e738d"),
		Overlay1:  lipgloss.Color("#8087a2"),
		Overlay2:  lipgloss.Color("#939ab7"),
		Subtext0:  lipgloss.Color("#a5adcb"),
		Subtext1:  lipgloss.Color("#b8c0e0"),
		Text:      lipgloss.Color("#cad3f5"),
		Rosewater: lipgloss.Color("#f4dbd6"),
		Flamingo:  lipgloss.Color("#f0c6c6"),
		Pink:      lipgloss.Color("#f5bde6"),
		Mauve:     lipgloss.Color("#c6a0f6"),
		Red:       lipgloss.Color("#ed8796"),
		Maroon:    lipgloss.Color("#ee99a0"),
		Peach:     lipgloss.Color("#f5a97f"),
		Yellow:    lipgloss.Color("#eed49f"),
		Green:     lipgloss.Color("#a6da95"),
		Teal:      lipgloss.Color("#8bd5ca"),
		Sky:       lipgloss.Color("#91d7e3"),
		Sapphire:  lipgloss.Color("#7dc4e4"),
		Blue:      lipgloss.Color("#8aadf4"),
		Lavender:  lipgloss.Color("#b7bdf8"),
	}
}

// Global theme instance - can be swapped later for theme switching
var Current = NewCatppuccinMacchiato()
