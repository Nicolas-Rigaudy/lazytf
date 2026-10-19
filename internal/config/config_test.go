package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadFrom(t *testing.T) {
	defaults := DefaultConfig()

	tests := []struct {
		name    string
		content *string // nil means no file
		want    Config
		wantErr bool
	}{
		{
			name: "missing file gives defaults",
			want: defaults,
		},
		{
			name:    "empty file gives defaults",
			content: ptr(""),
			want:    defaults,
		},
		{
			name:    "only search_paths keeps default ignore patterns",
			content: ptr("search_paths:\n  - ~/infra\n"),
			want:    Config{SearchPaths: []string{"~/infra"}, IgnorePatterns: defaults.IgnorePatterns},
		},
		{
			name:    "only ignore_patterns keeps default search paths",
			content: ptr("ignore_patterns:\n  - \"*/build\"\n"),
			want:    Config{SearchPaths: defaults.SearchPaths, IgnorePatterns: []string{"*/build"}},
		},
		{
			name:    "explicit empty ignore_patterns disables ignoring",
			content: ptr("search_paths: [\".\"]\nignore_patterns: []\n"),
			want:    Config{SearchPaths: []string{"."}, IgnorePatterns: []string{}},
		},
		{
			name:    "invalid yaml returns defaults and an error",
			content: ptr("search_paths: [unclosed\n"),
			want:    defaults,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yml")
			if tt.content != nil {
				if err := os.WriteFile(path, []byte(*tt.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			got, err := loadFrom(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !slices.Equal(got.SearchPaths, tt.want.SearchPaths) {
				t.Errorf("SearchPaths = %v, want %v", got.SearchPaths, tt.want.SearchPaths)
			}
			if !slices.Equal(got.IgnorePatterns, tt.want.IgnorePatterns) {
				t.Errorf("IgnorePatterns = %v, want %v", got.IgnorePatterns, tt.want.IgnorePatterns)
			}
		})
	}
}

func TestLoadFromDoesNotShareDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("search_paths: [\"/a\"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg.IgnorePatterns[0] = "mutated"
	if DefaultConfig().IgnorePatterns[0] == "mutated" {
		t.Error("mutating a loaded config changed DefaultConfig()")
	}
}

func ptr(s string) *string { return &s }
