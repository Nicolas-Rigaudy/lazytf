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
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
)

// DiscoverBackendVarFiles lists the backend config files matching
// layout.BackendConfigs. A missing project directory gives no files.
func DiscoverBackendVarFiles(projectPath string, layout config.Layout) ([]BackendVarFile, error) {
	var backendFiles []BackendVarFile
	envPattern, err := compileEnvPattern(layout.BackendEnvPattern)
	if err != nil {
		return nil, err
	}

	for _, rel := range globFiles(projectPath, layout.BackendConfigs) {
		name := path.Base(rel)
		relDir := filepath.FromSlash(path.Dir(rel))
		backendFiles = append(backendFiles, BackendVarFile{
			Name:     name,
			Path:     relDir,
			FullPath: filepath.Join(projectPath, filepath.FromSlash(rel)),
			EnvName:  extractBackendEnv(name, envPattern),
		})
	}
	return backendFiles, nil
}

// extractBackendEnv extracts the environment name from a backend config
// filename. With no custom pattern it follows the "backend_<env>" convention.
// Examples:
//   - "backend_dev2.tfvars" -> "dev2"
//   - "dev.hcl" -> "dev"
//   - "backend.tfvars" -> "" (generic, applies to all)
func extractBackendEnv(filename string, re *regexp.Regexp) string {
	if re != nil {
		return envName(filename, re)
	}
	return extractEnvFromBackendFile(filename)
}

// extractEnvFromBackendFile is the default backend naming convention.
func extractEnvFromBackendFile(filename string) string {
	name := stripVarExt(filename)

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
