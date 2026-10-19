package terraform

import (
	"slices"
	"testing"
)

func TestInferEnvFromBackendConfig(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]string
		want   string
	}{
		{name: "s3 key env dir", config: map[string]string{"key": "dev2/terraform.tfstate"}, want: "dev2"},
		{name: "s3 key env file", config: map[string]string{"key": "int.tfstate"}, want: "int"},
		{name: "workspace dir", config: map[string]string{"key": "terraform.tfstate.d/dev2"}, want: "dev2"},
		{name: "local path under states", config: map[string]string{"path": "states/int/terraform.tfstate"}, want: "int"},
		{name: "gcs prefix", config: map[string]string{"prefix": "prod"}, want: "prod"},
		{name: "workspace_key_prefix", config: map[string]string{"workspace_key_prefix": "qa"}, want: "qa"},
		{name: "leading slash", config: map[string]string{"key": "/dev/terraform.tfstate"}, want: "dev"},
		{name: "key wins over path", config: map[string]string{"key": "dev/terraform.tfstate", "path": "prod.tfstate"}, want: "dev"},
		{name: "generic state file only", config: map[string]string{"key": "terraform.tfstate"}, want: ""},
		{name: "unrelated keys", config: map[string]string{"bucket": "my-bucket", "region": "eu-west-1"}, want: ""},
		{name: "empty config", config: map[string]string{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferEnvFromBackendConfig(tt.config); got != tt.want {
				t.Errorf("inferEnvFromBackendConfig(%v) = %q, want %q", tt.config, got, tt.want)
			}
		})
	}
}

func TestDetectCurrentBackend(t *testing.T) {
	backends := []BackendVarFile{
		{Name: "backend_dev.tfvars", EnvName: "dev"},
		{Name: "backend_prod.tfvars", EnvName: "prod"},
	}

	tests := []struct {
		name        string
		files       map[string]string
		wantInit    bool
		wantType    string
		wantEnv     string
		wantMatched []string
	}{
		{name: "not initialised", files: map[string]string{"main.tf": ""}},
		{name: "initialised, no backend state", files: map[string]string{".terraform/providers/x": ""}, wantInit: true},
		{name: "initialised, invalid state json", files: map[string]string{".terraform/terraform.tfstate": "{"}, wantInit: true},
		{
			name: "s3 backend for dev",
			files: map[string]string{".terraform/terraform.tfstate": `{"backend":{"type":"s3","config":{
				"bucket":"b","key":"dev/terraform.tfstate","encrypt":true}}}`},
			wantInit: true, wantType: "s3", wantEnv: "dev", wantMatched: []string{"backend_dev.tfvars"},
		},
		{
			name:     "local backend, no env in path",
			files:    map[string]string{".terraform/terraform.tfstate": `{"backend":{"type":"local","config":{"path":"terraform.tfstate"}}}`},
			wantInit: true, wantType: "local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFiles(t, root, tt.files)
			got := DetectCurrentBackend(root, backends)
			if got.IsInitialized != tt.wantInit || got.BackendType != tt.wantType || got.DetectedEnv != tt.wantEnv {
				t.Errorf("DetectCurrentBackend() = init %v type %q env %q; want init %v type %q env %q",
					got.IsInitialized, got.BackendType, got.DetectedEnv, tt.wantInit, tt.wantType, tt.wantEnv)
			}
			var matched []string
			for _, b := range got.MatchedBackends {
				matched = append(matched, b.Name)
			}
			if !slices.Equal(matched, tt.wantMatched) {
				t.Errorf("MatchedBackends = %v, want %v", matched, tt.wantMatched)
			}
		})
	}
}
