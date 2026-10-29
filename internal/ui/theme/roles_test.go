package theme

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestActiveUsesOnlySixteenSlots(t *testing.T) {
	v := reflect.ValueOf(Active)
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		switch f := v.Field(i).Interface().(type) {
		case lipgloss.Color:
			checkSlot(t, name, f)
		case lipgloss.Style:
			if bg := f.GetBackground(); bg != (lipgloss.NoColor{}) {
				t.Errorf("%s paints a background %v", name, bg)
			}
			if fg, ok := f.GetForeground().(lipgloss.Color); ok {
				checkSlot(t, name, fg)
			}
		}
	}
}

func checkSlot(t *testing.T, name string, c lipgloss.Color) {
	t.Helper()
	if c == "" {
		return
	}
	if n, err := strconv.Atoi(string(c)); err != nil || n < 0 || n > 15 {
		t.Errorf("%s = %q, want an ANSI slot 0-15", name, c)
	}
}
