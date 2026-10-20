// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and resources.
//
// This file contains all project discovery and detection logic:
// - IsInitialized: Check if a directory has .terraform/
// - IsTerraformProject: Check if a directory contains .tf files
// - DetermineMode: Decide between single-project vs multi-project mode
// - DiscoverProjects: Find all Terraform projects in search paths
package terraform

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
)

// IsInitialized checks if a directory is an initialized Terraform project.
// A project is considered initialized if it contains a .terraform/ directory.
func IsInitialized(path string) bool {
	terraformDir := filepath.Join(path, ".terraform")
	info, err := os.Stat(terraformDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsTerraformProject checks if a directory contains Terraform files (*.tf or
// *.tf.json, the JSON syntax that generators such as CDKTF emit).
// This is more permissive than IsInitialized - it finds uninitialized projects too.
func IsTerraformProject(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && (strings.HasSuffix(name, ".tf") || strings.HasSuffix(name, ".tf.json")) {
			return true, nil
		}
	}
	return false, nil
}

// DetermineMode checks the current directory and decides which mode to run in.
// Returns ModeSingleProject if current directory is an initialized TF project,
// otherwise returns ModeMultiProject.
// If ModeSingleProject, also returns a Project struct with basic info.
func DetermineMode() (Mode, *Project, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return ModeMultiProject, nil, err
	}

	if IsInitialized(workDir) {
		project := Project{
			Name:             filepath.Base(workDir),
			Path:             workDir,
			IsInitialized:    true,
			Workspaces:       nil,
			CurrentWorkspace: "",
		}
		return ModeSingleProject, &project, nil
	}

	return ModeMultiProject, nil, nil
}

// ExpandPath expands a path that may contain ~ to the full home directory path.
// Example: ~/Projects -> /home/username/Projects
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	expandedPath := strings.Replace(path, "~", homeDir, 1)
	return expandedPath, nil
}

// shouldIgnore checks if a path matches any of the ignore patterns.
// A pattern is matched against as many trailing path segments as it has, so
// "*/.git" matches any .git directory and ".git" matches it too.
func shouldIgnore(path string, ignorePatterns []string) bool {
	segments := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")
	for _, pattern := range ignorePatterns {
		n := strings.Count(pattern, "/") + 1
		if n > len(segments) {
			continue
		}
		tail := strings.Join(segments[len(segments)-n:], "/")
		matched, err := filepath.Match(pattern, tail)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// DiscoverProjects scans the given search paths for Terraform projects.
// It respects the ignore patterns from the config.
func DiscoverProjects(cfg config.Config) ([]Project, error) {
	var projects []Project

	for _, searchPath := range cfg.SearchPaths {
		searchPath, err := ExpandPath(searchPath)
		if err != nil {
			return nil, err
		}
		err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip paths we can't read
			}

			// Check if we should ignore this path
			if shouldIgnore(path, cfg.IgnorePatterns) {
				if info.IsDir() {
					return filepath.SkipDir // Don't descend into ignored dirs
				}
				return nil // Skip ignored files
			}

			// Only check directories
			if !info.IsDir() {
				return nil
			}

			// Check if this directory is a TF project
			isTF, err := IsTerraformProject(path)
			if err != nil {
				return nil // Skip on error
			}

			if isTF {
				// Create and add project
				project := Project{
					Name:             filepath.Base(path),
					Path:             path,
					Workspaces:       nil,
					CurrentWorkspace: "",
					IsInitialized:    IsInitialized(path),
				}
				projects = append(projects, project)

				// Don't search inside TF projects (they're already found)
				return filepath.SkipDir
			}

			return nil // Continue walking
		})
		if err != nil {
			return nil, err
		}
	}

	// Remove duplicates and make project names unique
	projects = makeProjectNamesUnique(projects)

	return projects, nil
}

// makeProjectNamesUnique removes duplicate projects (by path) and ensures
// remaining project names are unique by prepending parent directory when needed.
func makeProjectNamesUnique(projects []Project) []Project {
	// First, deduplicate by path (remove actual duplicates)
	seen := make(map[string]bool)
	uniqueProjects := []Project{}

	for _, project := range projects {
		// Use absolute path for deduplication
		absPath := project.Path
		if !filepath.IsAbs(absPath) {
			// Convert relative paths to absolute for comparison
			if abs, err := filepath.Abs(absPath); err == nil {
				absPath = abs
			}
		}

		if !seen[absPath] {
			seen[absPath] = true
			uniqueProjects = append(uniqueProjects, project)
		}
	}

	// Now make names unique by prepending parent directory when needed
	nameCounts := make(map[string]int)
	for i := range uniqueProjects {
		nameCounts[uniqueProjects[i].Name]++
	}

	for i := range uniqueProjects {
		if nameCounts[uniqueProjects[i].Name] > 1 {
			parentDir := filepath.Base(filepath.Dir(uniqueProjects[i].Path))
			uniqueProjects[i].Name = parentDir + "/" + uniqueProjects[i].Name
		}
	}

	return uniqueProjects
}
