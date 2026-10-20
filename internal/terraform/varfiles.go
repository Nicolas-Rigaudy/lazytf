// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and resources.
//
// This file contains var file discovery logic:
// - DiscoverVarFiles: Find var files matching the project's layout
// - GetVarFileDisplayNames: Extract environment names for UI display
package terraform

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
)

// DiscoverVarFiles lists the environments of a project: the files matching
// layout.VarFiles, minus layout.Exclude, plus layout.Include. A missing
// project directory gives no files.
func DiscoverVarFiles(projectPath string, layout config.Layout) ([]VarFile, error) {
	var varFiles []VarFile
	envPattern, err := compileEnvPattern(layout.EnvPattern)
	if err != nil {
		return nil, err
	}

	include := map[string]bool{}
	for _, rel := range globFiles(projectPath, layout.Include) {
		include[rel] = true
	}
	candidates := append(globFiles(projectPath, layout.VarFiles), globFiles(projectPath, layout.Include)...)

	seen := map[string]bool{}
	for _, rel := range candidates {
		if seen[rel] {
			continue
		}
		seen[rel] = true
		if !include[rel] && excluded(layout.Exclude, rel) {
			continue
		}
		name := path.Base(rel)
		varFiles = append(varFiles, VarFile{
			Name:     name,
			Path:     filepath.FromSlash(rel),
			FullPath: filepath.Join(projectPath, filepath.FromSlash(rel)),
			EnvName:  envName(name, envPattern),
		})
	}
	return varFiles, nil
}

func excluded(patterns []string, rel string) bool {
	for _, p := range patterns {
		if matchExclude(p, rel) {
			return true
		}
	}
	return false
}

// compileEnvPattern compiles a user env pattern; empty gives nil (whole name).
func compileEnvPattern(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid env pattern %q: %w", pattern, err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("env pattern %q needs one capture group", pattern)
	}
	return re, nil
}

// stripVarExt removes the extension of a var or backend config file name:
// "dev.tfvars" -> "dev", "dev.tfvars.json" -> "dev", "dev.hcl" -> "dev".
func stripVarExt(name string) string {
	name = strings.TrimSuffix(name, ".json")
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// envName derives the environment name from a file name: the first capture
// group of re when it matches, otherwise the name without extension.
func envName(fileName string, re *regexp.Regexp) string {
	base := stripVarExt(fileName)
	if re != nil {
		if m := re.FindStringSubmatch(base); m != nil {
			return m[1]
		}
	}
	return base
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

func FindVarFileByEnvName(envName string, varFiles []VarFile) (*VarFile, int) {
	for i := range varFiles {
		if varFiles[i].EnvName == envName {
			return &varFiles[i], i
		}
	}
	return nil, -1
}
