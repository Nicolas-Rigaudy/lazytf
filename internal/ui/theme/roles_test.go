package theme

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestForNameDefaultsToCatppuccin(t *testing.T) {
	for _, n := range []string{"", "nope", "Catppuccin"} {
		if got := ForName(n).Name; got != "catppuccin" {
			t.Errorf("ForName(%q) = %s", n, got)
		}
	}
}

func TestANSIStylesUseOnlySixteenSlots(t *testing.T) {
	for _, n := range []string{"ansi-a", "ansi-b"} {
		r := ForName(n)
		if r.Name != n {
			t.Fatalf("ForName(%q).Name = %q", n, r.Name)
		}
		v := reflect.ValueOf(r)
		for i := 0; i < v.NumField(); i++ {
			if c, ok := v.Field(i).Interface().(lipgloss.Color); ok && c != "" {
				if num, err := strconv.Atoi(string(c)); err != nil || num < 0 || num > 15 {
					t.Errorf("%s.%s = %q, want an ANSI slot 0-15", n, v.Type().Field(i).Name, c)
				}
			}
		}
		for name, s := range map[string]lipgloss.Style{"Selected": r.Selected, "SelectedItem": r.SelectedItem, "TitleBar": r.TitleBar, "StatusBar": r.StatusBar, "ConfirmPrompt": r.ConfirmPrompt} {
			if bg := s.GetBackground(); bg != (lipgloss.NoColor{}) {
				t.Errorf("%s.%s paints a background %v", n, name, bg)
			}
			if fg, ok := s.GetForeground().(lipgloss.Color); ok && strings.HasPrefix(string(fg), "#") {
				t.Errorf("%s.%s uses hex %s", n, name, fg)
			}
		}
	}
}
