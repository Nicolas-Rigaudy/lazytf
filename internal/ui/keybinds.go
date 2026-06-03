package ui

type KeyBinding struct {
	Key         string
	Description string
	Enabled     bool
}

type KeyBindingGroup struct {
	Title       string
	KeyBindings []KeyBinding
}

func DefaultKeyBindings() []KeyBindingGroup {
	return []KeyBindingGroup{
		{
			Title: "Global",
			KeyBindings: []KeyBinding{
				{Key: "q", Description: "Quit", Enabled: true},
				{Key: "ctrl+c", Description: "Exit", Enabled: true},
				{Key: "backspace/esc", Description: "Back", Enabled: true},
				{Key: "?", Description: "Help", Enabled: true},
			},
		},
		{
			Title: "Navigation",
			KeyBindings: []KeyBinding{
				{Key: "h/left arrow", Description: "Scroll Left", Enabled: true},
				{Key: "l/right arrow", Description: "Scroll Right", Enabled: true},
				{Key: "j/down arrow", Description: "Scroll Down", Enabled: true},
				{Key: "k/up arrow", Description: "Scroll Up", Enabled: true},
				{Key: "enter", Description: "Select", Enabled: true},
				{Key: "tab", Description: "Switch Panel", Enabled: true},
				{Key: "b", Description: "Sidebar toggle", Enabled: true},
			},
		},
		{
			Title: "Terraform",
			KeyBindings: []KeyBinding{
				{Key: "L", Description: "Login", Enabled: true},
				{Key: "i", Description: "Init", Enabled: true},
				{Key: "v", Description: "Validate", Enabled: false},
				{Key: "f", Description: "Format", Enabled: false},
				{Key: "d", Description: "Destroy", Enabled: true},
				{Key: "p", Description: "Plan", Enabled: true},
				{Key: "a", Description: "Apply", Enabled: true},
			},
		},
	}
}
