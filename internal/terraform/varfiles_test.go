package terraform

import (
	"path/filepath"
	"slices"
	"testing"
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
			got, err := DiscoverVarFiles(root)
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
	got, err := DiscoverVarFiles(filepath.Join(t.TempDir(), "missing"))
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
