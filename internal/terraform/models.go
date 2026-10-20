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
	AnyEnv          bool              // true if no backend config files exist, so one init serves every env
}

// Mode represents how LazyTF is operating.
type Mode int

const (
	ModeSingleProject Mode = iota
	ModeMultiProject
)

// Binary returns the executable to run: the configured one, or "terraform".
// A name such as "tofu" is looked up on PATH.
func Binary(configured string) string {
	if configured == "" {
		return "terraform"
	}
	return configured
}

// NoVarFileEnv is the name of the stand-in environment of a project with no
// var files, so it can still be planned and applied.
const NoVarFileEnv = "(no var-file)"

// WithDefaultEnv returns varFiles, or one stand-in environment without a
// var file when there are none.
func WithDefaultEnv(varFiles []VarFile) []VarFile {
	if len(varFiles) > 0 {
		return varFiles
	}
	return []VarFile{{Name: NoVarFileEnv, EnvName: NoVarFileEnv}}
}
