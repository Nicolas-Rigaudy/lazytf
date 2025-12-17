// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and workspaces.
package terraform

// Project represents a Terraform project directory.
type Project struct {
	Name             string
	Path             string
	Workspaces       []string
	CurrentWorkspace string
	IsInitialized    bool
}

// VarFile represents a Terraform variables file (.tfvars).
type VarFile struct {
	Name     string // filename (e.g., "dev2.tfvars")
	Path     string // relative path from project root
	FullPath string // absolute path to the file
	EnvName  string // extracted environment name (e.g., "dev2")
}

// BackendVarFile represents a Terraform backend configuration file (.tfvars).
type BackendVarFile struct {
	Name     string // filename (e.g., "backend_dev2.tfvars")
	Path     string // relative path from project root (e.g., "variables/backend/local")
	FullPath string // absolute path to the file
	EnvName  string // extracted environment name (e.g., "dev2")
}

// BackendState represents the current Terraform backend initialization state.
type BackendState struct {
	IsInitialized   bool              // true if .terraform/ directory exists
	BackendType     string            // e.g., "s3", "local", "azurerm"
	BackendConfig   map[string]string // key-value pairs from backend config
	DetectedEnv     string            // environment name inferred from backend config (e.g., "dev2")
	MatchedBackends []BackendVarFile  // the backend var files that matches current state
}

// Mode represents how LazyTF is operating.
type Mode int

const (
	ModeSingleProject Mode = iota
	ModeMultiProject
)
