package terraform

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsLegacyPlanName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"lazytf-app.tfplan", true},
		{"lazytf-my-infra.tfplan", true},
		{"lazytf-.tfplan", true},
		{"lazytf-app.tfplan.bak", false},
		{"lazytf-app.plan", false},
		{"lazytf.tfplan", false},
		{"xlazytf-app.tfplan", false},
		{"app.tfplan", false},
		{"lazytf", false},
	}
	for _, tt := range tests {
		if got := isLegacyPlanName(tt.name); got != tt.want {
			t.Errorf("isLegacyPlanName(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("plan"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func TestRemoveLegacyPlans(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()

	plans := []string{"lazytf-app.tfplan", "lazytf-infra.tfplan"}
	for _, p := range plans {
		writeFile(t, filepath.Join(dir, p))
	}
	// Wrong pattern: kept.
	kept := []string{"lazytf-app.tfplan.bak", "other.tfplan", "lazytf-notes.txt"}
	for _, k := range kept {
		writeFile(t, filepath.Join(dir, k))
	}
	// Matching name in a subdirectory: no recursion, kept.
	sub := filepath.Join(dir, "lazytf", "plans")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(sub, "lazytf-app.tfplan")
	writeFile(t, nested)
	// A directory matching the pattern: not a regular file, kept.
	matchingDir := filepath.Join(dir, "lazytf-dir.tfplan")
	if err := os.Mkdir(matchingDir, 0o700); err != nil {
		t.Fatal(err)
	}
	// A symlink matching the pattern: neither the link nor its target goes.
	target := filepath.Join(outside, "lazytf-target.tfplan")
	writeFile(t, target)
	link := filepath.Join(dir, "lazytf-link.tfplan")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if got := removeLegacyPlans(dir, os.Getuid()); got != len(plans) {
		t.Errorf("removeLegacyPlans() = %d, want %d", got, len(plans))
	}
	for _, p := range plans {
		if exists(filepath.Join(dir, p)) {
			t.Errorf("%s still exists, want removed", p)
		}
	}
	for _, k := range kept {
		if !exists(filepath.Join(dir, k)) {
			t.Errorf("%s was removed, want kept (wrong pattern)", k)
		}
	}
	for _, p := range []string{nested, matchingDir, link, target} {
		if !exists(p) {
			t.Errorf("%s was removed, want kept", p)
		}
	}
}

func TestRemoveLegacyPlansSkipsOtherOwner(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no owner check on this platform")
	}
	dir := t.TempDir()
	plan := filepath.Join(dir, "lazytf-app.tfplan")
	writeFile(t, plan)

	// The file belongs to the current user; claim to be someone else.
	if got := removeLegacyPlans(dir, os.Getuid()+1); got != 0 {
		t.Errorf("removeLegacyPlans() = %d, want 0", got)
	}
	if !exists(plan) {
		t.Error("another user's plan file was removed")
	}
}

func TestRemoveLegacyPlansMissingDir(t *testing.T) {
	if got := removeLegacyPlans(filepath.Join(t.TempDir(), "missing"), os.Getuid()); got != 0 {
		t.Errorf("removeLegacyPlans() = %d, want 0", got)
	}
}

func TestLegacyPlansNote(t *testing.T) {
	tests := []struct {
		removed int
		want    string
	}{
		{0, ""},
		{1, "Removed 1 old plan file from /tmp (v0.1 location)"},
		{2, "Removed 2 old plan files from /tmp (v0.1 location)"},
	}
	for _, tt := range tests {
		if got := legacyPlansNote(tt.removed, "/tmp"); got != tt.want {
			t.Errorf("legacyPlansNote(%d) = %q, want %q", tt.removed, got, tt.want)
		}
	}
}
