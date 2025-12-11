// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and resources.
//
// This file contains var file discovery logic:
// - DiscoverVarFiles: Find .tfvars files in conventional locations
// - GetVarFileDisplayNames: Extract environment names for UI display
package terraform

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverVarFiles scans a Terraform project directory for .tfvars files
// in conventional locations (root, variables/, env/, tfvars/).
// Returns a list of discovered variable files.
func DiscoverVarFiles(projectPath string) ([]VarFile, error) {
	var varFiles []VarFile

	scanDirs := []string{".", "variables", "env", "tfvars"}

	for _, dir := range scanDirs {
		fullDirPath := filepath.Join(projectPath, dir)

		info, err := os.Stat(fullDirPath)
		if err != nil {
			// Directory doesn't exist, skip it
			continue
		}
		if !info.IsDir() {
			continue
		}

		entries, err := os.ReadDir(fullDirPath)
		if err != nil {
			// Can't read directory, skip it
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if filepath.Ext(entry.Name()) != ".tfvars" {
				continue
			}

			var relativePath string
			if dir == "." {
				relativePath = entry.Name()
			} else {
				relativePath = filepath.Join(dir, entry.Name())
			}

			// Extract environment name from filename
			// e.g., "dev2.tfvars" -> "dev2"
			envName := strings.TrimSuffix(entry.Name(), ".tfvars")

			varFile := VarFile{
				Name:     entry.Name(),
				Path:     relativePath,
				FullPath: filepath.Join(projectPath, relativePath),
				EnvName:  envName,
			}

			varFiles = append(varFiles, varFile)
		}
	}

	return varFiles, nil
}

// GetVarFileDisplayNames extracts environment names from var files for UI display.
// This helper keeps UI code clean and maintains separation of concerns.
func GetVarFileDisplayNames(varFiles []VarFile) []string {
	names := make([]string, len(varFiles))
	for i, vf := range varFiles {
		names[i] = vf.EnvName
	}
	return names
}
