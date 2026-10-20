package terraform

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// skipDirs are never searched for layout files: they hold downloaded
// modules, caches or VCS data, never a project's own var files.
var skipDirs = map[string]bool{".terraform": true, ".git": true, "node_modules": true, ".terragrunt-cache": true}

// matchGlob reports whether the slash-separated relative path matches pattern.
// It is path.Match per segment, plus a "**" segment that matches zero or more
// directories.
func matchGlob(pattern, rel string) bool {
	return matchSegments(strings.Split(pattern, "/"), strings.Split(rel, "/"))
}

func matchSegments(pat, segs []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			rest := pat[1:]
			for i := 0; i <= len(segs); i++ {
				if matchSegments(rest, segs[i:]) {
					return true
				}
			}
			return false
		}
		if len(segs) == 0 {
			return false
		}
		if ok, err := path.Match(pat[0], segs[0]); err != nil || !ok {
			return false
		}
		pat, segs = pat[1:], segs[1:]
	}
	return len(segs) == 0
}

// matchExclude matches a glob against a relative path: without a "/" it is
// compared with the file name alone, so "*.auto.tfvars" excludes it anywhere.
func matchExclude(pattern, rel string) bool {
	if !strings.Contains(pattern, "/") {
		return matchGlob(pattern, path.Base(rel))
	}
	return matchGlob(pattern, rel)
}

// globFiles returns the regular files under root matching any pattern, as
// slash-separated paths relative to root. Results follow pattern order, and
// are sorted within a pattern; a file matched twice is listed once.
func globFiles(root string, patterns []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, pattern := range patterns {
		pattern = path.Clean(filepath.ToSlash(pattern))
		start, recursive := globBase(pattern)
		var matches []string
		_ = filepath.WalkDir(filepath.Join(root, filepath.FromSlash(start)), func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // unreadable or missing: skip
			}
			if d.IsDir() {
				if skipDirs[d.Name()] && p != root {
					return filepath.SkipDir
				}
				if !recursive && p != filepath.Join(root, filepath.FromSlash(start)) {
					return filepath.SkipDir
				}
				return nil
			}
			if !d.Type().IsRegular() {
				// Symlinks to files are common for shared var files.
				if info, err := os.Stat(p); err != nil || !info.Mode().IsRegular() {
					return nil
				}
			}
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)
			if matchGlob(pattern, rel) {
				matches = append(matches, rel)
			}
			return nil
		})
		sort.Strings(matches)
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	return out
}

// globBase returns the directory to start walking from (the segments before
// the first wildcard) and whether the walk must descend past it.
func globBase(pattern string) (start string, recursive bool) {
	segs := strings.Split(pattern, "/")
	var fixed []string
	for i, s := range segs {
		if s == "**" {
			return path.Join(append(fixed, ".")...), true
		}
		if strings.ContainsAny(s, "*?[") {
			// A wildcard in a directory segment needs descent.
			return path.Join(append(fixed, ".")...), i < len(segs)-1
		}
		if i < len(segs)-1 {
			fixed = append(fixed, s)
		}
	}
	return path.Join(append(fixed, ".")...), false
}
