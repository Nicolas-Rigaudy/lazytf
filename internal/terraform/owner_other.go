//go:build !unix

package terraform

import "os"

// ownedBy has no portable owner check here; the temp dir is per-user.
func ownedBy(os.FileInfo, int) bool { return true }
