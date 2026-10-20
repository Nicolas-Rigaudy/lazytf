package ui

// KeyBinding is one entry of the help modal. The status bar's hints come from
// the same list (Hint), so the two cannot drift apart.
type KeyBinding struct {
	Key         string
	Description string
	Enabled     bool
	Hint        string // status bar text in detail view, "" = not shown there
}

type KeyBindingGroup struct {
	Title       string
	KeyBindings []KeyBinding
}

// DefaultKeyBindings lists the help-modal keys; the login key is listed only
// when an auth provider is configured.
func DefaultKeyBindings(loginEnabled bool) []KeyBindingGroup {
	groups := []KeyBindingGroup{
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
				{Key: "i", Description: "Init", Enabled: true, Hint: "i: init"},
				{Key: "v", Description: "Toggle plan summary / raw output", Enabled: true},
				{Key: "d", Description: "Destroy", Enabled: true, Hint: "d: destroy"},
				{Key: "p", Description: "Plan", Enabled: true, Hint: "p: plan"},
				{Key: "a", Description: "Apply", Enabled: true, Hint: "a: apply"},
			},
		},
	}
	if loginEnabled {
		terraform := &groups[2].KeyBindings
		*terraform = append([]KeyBinding{{Key: "L", Description: "Login", Enabled: true, Hint: "L: login"}}, *terraform...)
	}
	return groups
}

// detailHints returns the status bar hints for the Terraform keys, in the
// order they are listed in the help modal.
func detailHints(loginEnabled bool) []string {
	var hints []string
	for _, group := range DefaultKeyBindings(loginEnabled) {
		for _, kb := range group.KeyBindings {
			if kb.Enabled && kb.Hint != "" {
				hints = append(hints, kb.Hint)
			}
		}
	}
	return hints
}
