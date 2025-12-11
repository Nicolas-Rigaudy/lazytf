// Package terraform provides functionality for detecting, discovering,
// and interacting with Terraform projects and resources.
//
// This file contains Terraform state detection logic:
// - DetectCurrentBackend: Read .terraform/terraform.tfstate to detect current backend
// - inferEnvFromBackendConfig: Extract environment name from backend config
// - FormatBackendState: Format backend state for UI display
package terraform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// DetectCurrentBackend checks if Terraform is initialized and detects the current backend configuration.
// It reads .terraform/terraform.tfstate to determine which backend is currently active.
func DetectCurrentBackend(projectPath string, backendVarFiles []BackendVarFile) BackendState {
	state := BackendState{
		IsInitialized: false,
		BackendConfig: make(map[string]string),
	}

	// Check if .terraform directory exists
	terraformDir := filepath.Join(projectPath, ".terraform")
	if _, err := os.Stat(terraformDir); os.IsNotExist(err) {
		// Not initialized
		return state
	}

	// Directory exists - project is initialized
	state.IsInitialized = true

	// Try to read .terraform/terraform.tfstate
	tfstatePath := filepath.Join(terraformDir, "terraform.tfstate")
	data, err := os.ReadFile(tfstatePath)
	if err != nil {
		// File doesn't exist or can't be read - still initialized, just no backend info
		return state
	}

	// Parse the JSON structure
	var tfstate struct {
		Backend struct {
			Type   string                 `json:"type"`
			Config map[string]interface{} `json:"config"`
		} `json:"backend"`
	}

	if err := json.Unmarshal(data, &tfstate); err != nil {
		// Invalid JSON, return what we have
		return state
	}

	// Extract backend type
	state.BackendType = tfstate.Backend.Type

	// Convert config to string map for easier handling
	for key, value := range tfstate.Backend.Config {
		if strValue, ok := value.(string); ok {
			state.BackendConfig[key] = strValue
		}
	}

	// Try to infer environment name from backend config
	// Common patterns:
	// - S3 key: "dev2/terraform.tfstate" or "terraform.tfstate.d/dev2"
	// - File path: "terraform-dev2.tfstate"
	state.DetectedEnv = inferEnvFromBackendConfig(state.BackendConfig)

	// Try to match with available backend var files
	if state.DetectedEnv != "" {
		for i := range backendVarFiles {
			if backendVarFiles[i].EnvName == state.DetectedEnv {
				state.MatchedBackend = &backendVarFiles[i]
				break
			}
		}
	}

	return state
}

// inferEnvFromBackendConfig tries to extract environment name from backend configuration.
// Examples:
//   - S3 key "dev2/terraform.tfstate" -> "dev2"
//   - S3 key "int.tfstate" -> "int"
//   - S3 key "terraform.tfstate.d/dev2" -> "dev2"
//   - Path "states/int/terraform.tfstate" -> "int"
func inferEnvFromBackendConfig(config map[string]string) string {
	// Check common keys that might contain env info
	candidates := []string{"key", "path", "prefix", "workspace_key_prefix"}

	for _, key := range candidates {
		if value, exists := config[key]; exists {
			// Try to extract env from path segments
			parts := strings.Split(value, "/")
			for _, part := range parts {
				// Skip empty parts
				if part == "" {
					continue
				}

				// Strip .tfstate suffix if present (handles "int.tfstate" -> "int")
				cleanPart := strings.TrimSuffix(part, ".tfstate")

				// Skip generic terraform state file
				if cleanPart == "terraform" {
					continue
				}

				// Skip .d suffix directory names
				if strings.HasSuffix(cleanPart, ".d") {
					continue
				}

				// Skip common non-env directory names
				if cleanPart == "states" || cleanPart == "terraform" {
					continue
				}

				// This looks like an env name - return it
				// Common patterns: "dev2", "int", "prod"
				if len(cleanPart) > 0 {
					return cleanPart
				}
			}
		}
	}

	return ""
}

// FormatBackendState creates a user-friendly string describing the current backend state.
// Shows initialization status, backend type, and which environment is currently active.
func FormatBackendState(state BackendState) string {
	if !state.IsInitialized {
		return "❌ Not initialized\n\nRun 'terraform init' with a backend config to get started."
	}

	result := "✅ Initialized\n\n"

	if state.BackendType != "" {
		result += "Backend Type: " + state.BackendType + "\n"
	}

	if state.DetectedEnv != "" {
		result += "Current Environment: " + state.DetectedEnv + "\n"
	}

	if state.MatchedBackend != nil {
		result += "Backend Config: " + state.MatchedBackend.Name + "\n"
		result += "Config Path: " + state.MatchedBackend.Path + "\n"
	}

	return result
}
