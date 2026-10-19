package terraform

import (
	"fmt"
	"os"
	"path/filepath"
)

// legacyPlanPattern matches the plan files v0.1 saved in the temp dir
// (/tmp/lazytf-<dir>.tfplan). Plans now live under the user cache dir.
const legacyPlanPattern = "lazytf-*.tfplan"

// CleanupLegacyPlans deletes v0.1 plan files left in os.TempDir() and
// returns a status-bar note, or "" when nothing was removed.
// Best effort: every error is ignored.
func CleanupLegacyPlans() string {
	dir := os.TempDir()
	return legacyPlansNote(removeLegacyPlans(dir, os.Getuid()), dir)
}

func legacyPlansNote(removed int, dir string) string {
	switch removed {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("Removed 1 old plan file from %s (v0.1 location)", dir)
	default:
		return fmt.Sprintf("Removed %d old plan files from %s (v0.1 location)", removed, dir)
	}
}

// removeLegacyPlans deletes, directly in dir (no recursion), the regular
// files matching legacyPlanPattern that are owned by uid. Symlinks,
// directories and other users' files are left alone.
func removeLegacyPlans(dir string, uid int) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	removed := 0
	for _, e := range entries {
		if !isLegacyPlanName(e.Name()) {
			continue
		}
		// Info is an lstat: a symlink reports as a symlink, not its target.
		info, err := e.Info()
		if err != nil || !info.Mode().IsRegular() || !ownedBy(info, uid) {
			continue
		}
		if os.Remove(filepath.Join(dir, e.Name())) == nil {
			removed++
		}
	}
	return removed
}

func isLegacyPlanName(name string) bool {
	ok, _ := filepath.Match(legacyPlanPattern, name)
	return ok
}
