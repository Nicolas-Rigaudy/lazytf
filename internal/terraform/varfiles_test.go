package terraform

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
)

func TestDiscoverVarFiles(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  []VarFile // FullPath is checked as root + Path
	}{
		{
			name:  "no var files",
			files: map[string]string{"main.tf": ""},
			want:  nil,
		},
		{
			name: "all conventional dirs, in scan order",
			files: map[string]string{
				"root.tfvars":           "",
				"variables/dev.tfvars":  "",
				"variables/prod.tfvars": "",
				"env/int.tfvars":        "",
				"tfvars/qa.tfvars":      "",
			},
			want: []VarFile{
				{Name: "root.tfvars", Path: "root.tfvars", EnvName: "root"},
				{Name: "dev.tfvars", Path: "variables/dev.tfvars", EnvName: "dev"},
				{Name: "prod.tfvars", Path: "variables/prod.tfvars", EnvName: "prod"},
				{Name: "int.tfvars", Path: "env/int.tfvars", EnvName: "int"},
				{Name: "qa.tfvars", Path: "tfvars/qa.tfvars", EnvName: "qa"},
			},
		},
		{
			name: "ignores other files, nested dirs and unknown dirs",
			files: map[string]string{
				"variables/dev.tfvars":             "",
				"variables/dev.tfvars.json":        "",
				"variables/notes.md":               "",
				"variables/backend/local/x.tfvars": "", // backend configs are not var files
				"config/staging.tfvars":            "",
			},
			want: []VarFile{
				{Name: "dev.tfvars", Path: "variables/dev.tfvars", EnvName: "dev"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFiles(t, root, tt.files)
			got, err := DiscoverVarFiles(root, config.DefaultLayout())
			if err != nil {
				t.Fatal(err)
			}
			var want []VarFile
			for _, w := range tt.want {
				w.Path = filepath.FromSlash(w.Path)
				w.FullPath = filepath.Join(root, w.Path)
				want = append(want, w)
			}
			if !slices.Equal(got, want) {
				t.Errorf("DiscoverVarFiles() =\n  %+v\nwant\n  %+v", got, want)
			}
		})
	}
}

func TestDiscoverVarFilesMissingProject(t *testing.T) {
	got, err := DiscoverVarFiles(filepath.Join(t.TempDir(), "missing"), config.DefaultLayout())
	if err != nil || len(got) != 0 {
		t.Errorf("DiscoverVarFiles() = %v, %v; want none and no error", got, err)
	}
}

func TestFindVarFileByEnvName(t *testing.T) {
	files := []VarFile{{EnvName: "dev"}, {EnvName: "prod"}}
	if vf, i := FindVarFileByEnvName("prod", files); i != 1 || vf != &files[1] {
		t.Errorf("FindVarFileByEnvName(prod) = %v, %d; want &files[1], 1", vf, i)
	}
	if vf, i := FindVarFileByEnvName("qa", files); i != -1 || vf != nil {
		t.Errorf("FindVarFileByEnvName(qa) = %v, %d; want nil, -1", vf, i)
	}
}

func TestDiscoverVarFilesLayout(t *testing.T) {
	tests := []struct {
		name   string
		files  map[string]string
		layout config.Layout
		want   []string // env names, in order
	}{
		{
			name:   "default layout skips terraform.tfvars and auto.tfvars",
			files:  map[string]string{"terraform.tfvars": "", "x.auto.tfvars": "", "dev.tfvars": ""},
			layout: config.DefaultLayout(),
			want:   []string{"dev"},
		},
		{
			name:   "custom glob for a stranger's layout",
			files:  map[string]string{"envs/dev.tfvars": "", "envs/prod.tfvars": "", "dev.tfvars": ""},
			layout: config.Layout{VarFiles: []string{"envs/*.tfvars"}},
			want:   []string{"dev", "prod"},
		},
		{
			name:   "double star spans directories",
			files:  map[string]string{"a/b/qa.tfvars": "", "a/dev.tfvars": "", ".terraform/modules/x.tfvars": ""},
			layout: config.Layout{VarFiles: []string{"**/*.tfvars"}},
			want:   []string{"dev", "qa"},
		},
		{
			name:   "include wins over exclude",
			files:  map[string]string{"terraform.tfvars": "", "dev.tfvars": ""},
			layout: config.Layout{VarFiles: []string{"*.tfvars"}, Exclude: []string{"*.tfvars"}, Include: []string{"terraform.tfvars"}},
			want:   []string{"terraform"},
		},
		{
			name:   "exclude with a slash matches the relative path",
			files:  map[string]string{"envs/shared.tfvars": "", "envs/dev.tfvars": ""},
			layout: config.Layout{VarFiles: []string{"envs/*.tfvars"}, Exclude: []string{"envs/shared.tfvars"}},
			want:   []string{"dev"},
		},
		{
			name:   "env pattern captures the name",
			files:  map[string]string{"vars-dev.tfvars": "", "vars-prod.tfvars": "", "other.tfvars": ""},
			layout: config.Layout{VarFiles: []string{"*.tfvars"}, EnvPattern: `^vars-(.+)$`},
			want:   []string{"other", "dev", "prod"},
		},
		{
			name:   "tfvars.json names drop both extensions",
			files:  map[string]string{"dev.tfvars.json": ""},
			layout: config.Layout{VarFiles: []string{"*.tfvars.json"}},
			want:   []string{"dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFiles(t, root, tt.files)
			got, err := DiscoverVarFiles(root, tt.layout)
			if err != nil {
				t.Fatal(err)
			}
			names := GetVarFileDisplayNames(got)
			slices.Sort(names)
			want := slices.Clone(tt.want)
			slices.Sort(want)
			if !slices.Equal(names, want) {
				t.Errorf("envs = %v, want %v", names, want)
			}
		})
	}
}

func TestDiscoverVarFilesBadEnvPattern(t *testing.T) {
	for _, pattern := range []string{"(", "nogroup"} {
		_, err := DiscoverVarFiles(t.TempDir(), config.Layout{EnvPattern: pattern})
		if err == nil {
			t.Errorf("EnvPattern %q: want an error", pattern)
		}
	}
}
