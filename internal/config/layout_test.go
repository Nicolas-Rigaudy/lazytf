package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadLayout(t *testing.T) {
	base := DefaultLayout()
	tests := []struct {
		name    string
		content *string // nil means no file
		check   func(t *testing.T, got Layout)
		wantErr bool
	}{
		{
			name: "missing file gives the base",
			check: func(t *testing.T, got Layout) {
				if !slices.Equal(got.VarFiles, base.VarFiles) || !slices.Equal(got.Exclude, base.Exclude) {
					t.Errorf("got %+v, want base", got)
				}
			},
		},
		{
			name:    "three-line layout overrides only what it sets",
			content: ptr("var_files: [\"envs/*.tfvars\"]\nbackend_configs: [\"backends/*.hcl\"]\nenv_pattern: \"^(.+)$\"\n"),
			check: func(t *testing.T, got Layout) {
				if !slices.Equal(got.VarFiles, []string{"envs/*.tfvars"}) {
					t.Errorf("VarFiles = %v", got.VarFiles)
				}
				if !slices.Equal(got.BackendConfigs, []string{"backends/*.hcl"}) {
					t.Errorf("BackendConfigs = %v", got.BackendConfigs)
				}
				if !slices.Equal(got.Exclude, base.Exclude) {
					t.Errorf("Exclude = %v, want the base %v", got.Exclude, base.Exclude)
				}
			},
		},
		{
			name:    "explicit empty exclude clears the base",
			content: ptr("exclude: []\n"),
			check: func(t *testing.T, got Layout) {
				if got.Exclude == nil || len(got.Exclude) != 0 {
					t.Errorf("Exclude = %#v, want empty", got.Exclude)
				}
			},
		},
		{
			name:    "terraform binary can be set per project",
			content: ptr("terraform: tofu\n"),
			check: func(t *testing.T, got Layout) {
				if got.Terraform != "tofu" {
					t.Errorf("Terraform = %q, want tofu", got.Terraform)
				}
			},
		},
		{
			name:    "auth section",
			content: ptr("auth:\n  provider: gcloud\n  error_patterns: [\"corp expired\"]\n"),
			check: func(t *testing.T, got Layout) {
				if got.Auth.Provider != "gcloud" || !slices.Equal(got.Auth.ErrorPatterns, []string{"corp expired"}) {
					t.Errorf("Auth = %+v", got.Auth)
				}
			},
		},
		{
			name:    "invalid yaml returns the base and an error",
			content: ptr("var_files: [unclosed\n"),
			wantErr: true,
			check: func(t *testing.T, got Layout) {
				if !slices.Equal(got.VarFiles, base.VarFiles) {
					t.Errorf("VarFiles = %v, want base", got.VarFiles)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.content != nil {
				if err := os.WriteFile(filepath.Join(dir, ProjectFileName), []byte(*tt.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := LoadLayout(dir, base)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.check(t, got)
		})
	}
}

func TestLoadLayoutDoesNotShareBase(t *testing.T) {
	base := DefaultLayout()
	got, err := LoadLayout(t.TempDir(), base)
	if err != nil {
		t.Fatal(err)
	}
	got.VarFiles[0] = "mutated"
	if base.VarFiles[0] == "mutated" {
		t.Error("mutating a loaded layout changed the base")
	}
}

func TestGlobalConfigCarriesLayoutKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("var_files: [\"*.tfvars\"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(cfg.VarFiles, []string{"*.tfvars"}) {
		t.Errorf("VarFiles = %v", cfg.VarFiles)
	}
	if !slices.Equal(cfg.BackendConfigs, DefaultLayout().BackendConfigs) {
		t.Errorf("BackendConfigs = %v, want the default", cfg.BackendConfigs)
	}
}

func TestLoadLayoutDoesNotShareAuthEnv(t *testing.T) {
	base := DefaultLayout()
	base.Auth.Env = map[string]map[string]string{"dev": {"AWS_PROFILE": "a"}}
	got, err := LoadLayout(t.TempDir(), base)
	if err != nil {
		t.Fatal(err)
	}
	got.Auth.Env["dev"]["AWS_PROFILE"] = "mutated"
	if base.Auth.Env["dev"]["AWS_PROFILE"] != "a" {
		t.Error("mutating a loaded layout changed the base auth env")
	}
}
