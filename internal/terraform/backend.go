// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and resources.
//
// This file contains backend var file discovery and matching logic:
// - DiscoverBackendVarFiles: Find backend .tfvars files
// - extractEnvFromBackendFile: Extract environment name from filename
// - MatchBackendsForEnv: Match backend configs to environment names
// - FormatBackendInfo: Format backend info for UI display
package terraform

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverBackendVarFiles scans for backend configuration files in the project.
// It looks for .tfvars files in variables/backend/** subdirectories.
func DiscoverBackendVarFiles(projectPath string) ([]BackendVarFile, error) {
	var backendFiles []BackendVarFile

	// Look for backend configs in variables/backend/**/*.tfvars
	backendDir := filepath.Join(projectPath, "variables", "backend")

	// Check if backend directory exists
	if _, err := os.Stat(backendDir); os.IsNotExist(err) {
		// No backend directory, return empty list
		return backendFiles, nil
	}

	// Walk the backend directory recursively
	err := filepath.Walk(backendDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip paths we can't read
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .tfvars files
		if !strings.HasSuffix(info.Name(), ".tfvars") {
			return nil
		}

		// Get relative path from project root
		relPath, err := filepath.Rel(projectPath, filepath.Dir(path))
		if err != nil {
			relPath = filepath.Dir(path)
		}

		// Extract environment name from filename
		// e.g., "backend_dev2.tfvars" -> "dev2"
		// or "backend.tfvars" -> "" (generic backend)
		envName := extractEnvFromBackendFile(info.Name())

		backendFile := BackendVarFile{
			Name:     info.Name(),
			Path:     relPath,
			FullPath: path,
			EnvName:  envName,
		}

		backendFiles = append(backendFiles, backendFile)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return backendFiles, nil
}

// extractEnvFromBackendFile extracts the environment name from a backend filename.
// Examples:
//   - "backend_dev2.tfvars" -> "dev2"
//   - "backend_int.tfvars" -> "int"
//   - "backend.tfvars" -> "" (generic, applies to all)
func extractEnvFromBackendFile(filename string) string {
	// Remove .tfvars extension
	name := strings.TrimSuffix(filename, ".tfvars")

	// Remove "backend_" prefix if present
	if strings.HasPrefix(name, "backend_") {
		return strings.TrimPrefix(name, "backend_")
	}

	// If it's just "backend", it's a generic backend config
	if name == "backend" {
		return ""
	}

	// Otherwise return the name as-is
	return name
}

// MatchBackendsForEnv finds backend config files that match the given environment name.
// Matching logic:
//   - Find backend files with matching EnvName
//   - If no match found, return generic backend.tfvars if it exists
func MatchBackendsForEnv(envName string, backends []BackendVarFile) []BackendVarFile {
	var matches []BackendVarFile
	var genericBackend *BackendVarFile

	for i := range backends {
		backend := &backends[i]

		// Check if this is a generic backend (no env name)
		if backend.EnvName == "" {
			genericBackend = backend
			continue
		}

		// Check if env name matches
		if backend.EnvName == envName {
			matches = append(matches, *backend)
		}
	}

	// If we found matches, return them
	if len(matches) > 0 {
		return matches
	}

	// If no matches, return generic backend if it exists
	if genericBackend != nil {
		return []BackendVarFile{*genericBackend}
	}

	// No backend configs found
	return []BackendVarFile{}
}

// FormatBackendInfo creates a human-readable string describing backend configuration(s).
// Handles cases: no backends, single backend, multiple backends.
func FormatBackendInfo(backends []BackendVarFile) string {
	if len(backends) == 0 {
		return "No backend configuration found"
	}

	if len(backends) == 1 {
		b := backends[0]
		// Show name and location (e.g., "backend_dev2.tfvars (variables/backend/local)")
		return b.Name + " (" + b.Path + ")"
	}

	// Multiple backends found
	result := "Multiple backend options available:\n"
	for _, b := range backends {
		result += "  • " + b.Name + " (" + b.Path + ")\n"
	}
	return result
}
