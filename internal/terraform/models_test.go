package terraform

import "testing"

func TestBinary(t *testing.T) {
	for in, want := range map[string]string{"": "terraform", "tofu": "tofu", "/opt/tf/terraform": "/opt/tf/terraform"} {
		if got := Binary(in); got != want {
			t.Errorf("Binary(%q) = %q, want %q", in, got, want)
		}
	}
}
