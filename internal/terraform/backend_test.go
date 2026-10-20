package terraform

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
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
	got, err := DiscoverBackendVarFiles(root, config.DefaultLayout())
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
	got, err := DiscoverBackendVarFiles(t.TempDir(), config.DefaultLayout())
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

func TestDiscoverBackendVarFilesLayout(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"backends/dev.hcl":  "",
		"backends/prod.hcl": "",
		"backends/notes.md": "",
	})
	got, err := DiscoverBackendVarFiles(root, config.Layout{BackendConfigs: []string{"backends/*.hcl"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []BackendVarFile{
		{Name: "dev.hcl", Path: "backends", FullPath: filepath.Join(root, "backends", "dev.hcl"), EnvName: "dev"},
		{Name: "prod.hcl", Path: "backends", FullPath: filepath.Join(root, "backends", "prod.hcl"), EnvName: "prod"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestDiscoverBackendVarFilesEnvPattern(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{"be/state-dev.hcl": ""})
	got, err := DiscoverBackendVarFiles(root, config.Layout{BackendConfigs: []string{"be/*.hcl"}, BackendEnvPattern: `^state-(.+)$`})
	if err != nil || len(got) != 1 || got[0].EnvName != "dev" {
		t.Errorf("got %+v, %v; want one backend for env dev", got, err)
	}
}

func TestMatchGlob(t *testing.T) {
	tests := []struct {
		pattern, rel string
		want         bool
	}{
		{"*.tfvars", "dev.tfvars", true},
		{"*.tfvars", "a/dev.tfvars", false},
		{"**/*.tfvars", "dev.tfvars", true},
		{"**/*.tfvars", "a/b/dev.tfvars", true},
		{"a/**/x.hcl", "a/x.hcl", true},
		{"a/**/x.hcl", "a/b/c/x.hcl", true},
		{"a/**/x.hcl", "b/x.hcl", false},
		{"envs/[dp]*.tfvars", "envs/dev.tfvars", true},
		{"envs/[dp]*.tfvars", "envs/qa.tfvars", false},
	}
	for _, tt := range tests {
		if got := matchGlob(tt.pattern, tt.rel); got != tt.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", tt.pattern, tt.rel, got, tt.want)
		}
	}
}
