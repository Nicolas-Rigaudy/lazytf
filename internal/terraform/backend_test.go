package terraform

import (
	"path/filepath"
	"slices"
	"testing"
)

func TestExtractEnvFromBackendFile(t *testing.T) {
	tests := []struct{ in, want string }{
		{"backend_dev2.tfvars", "dev2"},
		{"backend_int.tfvars", "int"},
		{"backend.tfvars", ""},
		{"prod.tfvars", "prod"},
		{"backend_.tfvars", ""},
	}
	for _, tt := range tests {
		if got := extractEnvFromBackendFile(tt.in); got != tt.want {
			t.Errorf("extractEnvFromBackendFile(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDiscoverBackendVarFiles(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"variables/backend/local/backend_dev.tfvars": "",
		"variables/backend/s3/backend_dev.tfvars":    "",
		"variables/backend/backend.tfvars":           "",
		"variables/backend/s3/README.md":             "",
		"variables/dev.tfvars":                       "", // a var file, not a backend config
	})
	got, err := DiscoverBackendVarFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []BackendVarFile{
		{Name: "backend.tfvars", Path: filepath.FromSlash("variables/backend"), EnvName: ""},
		{Name: "backend_dev.tfvars", Path: filepath.FromSlash("variables/backend/local"), EnvName: "dev"},
		{Name: "backend_dev.tfvars", Path: filepath.FromSlash("variables/backend/s3"), EnvName: "dev"},
	}
	for i := range want {
		want[i].FullPath = filepath.Join(root, want[i].Path, want[i].Name)
	}
	if !slices.Equal(got, want) {
		t.Errorf("DiscoverBackendVarFiles() =\n  %+v\nwant\n  %+v", got, want)
	}
}

func TestDiscoverBackendVarFilesNoBackendDir(t *testing.T) {
	got, err := DiscoverBackendVarFiles(t.TempDir())
	if err != nil || len(got) != 0 {
		t.Errorf("DiscoverBackendVarFiles() = %v, %v; want none and no error", got, err)
	}
}

func TestMatchBackendsForEnv(t *testing.T) {
	devLocal := BackendVarFile{Name: "backend_dev.tfvars", Path: "variables/backend/local", EnvName: "dev"}
	devS3 := BackendVarFile{Name: "backend_dev.tfvars", Path: "variables/backend/s3", EnvName: "dev"}
	prod := BackendVarFile{Name: "backend_prod.tfvars", Path: "variables/backend/s3", EnvName: "prod"}
	generic := BackendVarFile{Name: "backend.tfvars", Path: "variables/backend", EnvName: ""}

	tests := []struct {
		name     string
		env      string
		backends []BackendVarFile
		want     []BackendVarFile
	}{
		{name: "single match", env: "prod", backends: []BackendVarFile{devLocal, prod, generic}, want: []BackendVarFile{prod}},
		{name: "several matches keep order", env: "dev", backends: []BackendVarFile{devLocal, prod, devS3}, want: []BackendVarFile{devLocal, devS3}},
		{name: "env match wins over generic", env: "dev", backends: []BackendVarFile{generic, devLocal}, want: []BackendVarFile{devLocal}},
		{name: "falls back to generic", env: "qa", backends: []BackendVarFile{devLocal, generic}, want: []BackendVarFile{generic}},
		{name: "no match and no generic", env: "qa", backends: []BackendVarFile{devLocal, prod}, want: []BackendVarFile{}},
		{name: "no backends", env: "dev", backends: nil, want: []BackendVarFile{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchBackendsForEnv(tt.env, tt.backends)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("MatchBackendsForEnv(%q) = %+v, want %+v", tt.env, got, tt.want)
			}
		})
	}
}
