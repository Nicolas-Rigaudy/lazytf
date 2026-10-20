package terraform

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Nicolas-Rigaudy/lazytf/internal/config"
)

func TestShouldIgnore(t *testing.T) {
	defaults := config.DefaultConfig().IgnorePatterns

	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{name: "default .git", path: "/home/u/repo/.git", patterns: defaults, want: true},
		{name: "default node_modules", path: "/home/u/app/node_modules", patterns: defaults, want: true},
		{name: "default .terraform", path: "/home/u/infra/.terraform", patterns: defaults, want: true},
		{name: "default vendor", path: "/home/u/go/vendor", patterns: defaults, want: true},
		{name: "relative path", path: "repo/.git", patterns: defaults, want: true},
		{name: "regular dir", path: "/home/u/infra/modules", patterns: defaults, want: false},
		{name: "name only contains pattern", path: "/home/u/.github", patterns: defaults, want: false},
		{name: "search root itself", path: ".", patterns: defaults, want: false},
		{name: "basename pattern", path: "/home/u/repo/.git", patterns: []string{".git"}, want: true},
		{name: "glob basename pattern", path: "/home/u/build-tmp", patterns: []string{"build-*"}, want: true},
		{name: "two-level pattern", path: "/home/u/examples/old", patterns: []string{"examples/*"}, want: true},
		{name: "two-level pattern, other parent", path: "/home/u/live/old", patterns: []string{"examples/*"}, want: false},
		{name: "no patterns", path: "/home/u/repo/.git", patterns: nil, want: false},
		{name: "malformed pattern is skipped", path: "/home/u/repo", patterns: []string{"["}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldIgnore(tt.path, tt.patterns); got != tt.want {
				t.Errorf("shouldIgnore(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestDiscoverProjects(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"live/network/main.tf":                    "",
		"live/network/modules/vpc/main.tf":        "", // inside a project: not a separate project
		"live/app/main.tf":                        "",
		"live/app/.terraform/terraform.tfstate":   "{}",
		"other/app/main.tf":                       "", // same basename as live/app
		"docs/README.md":                          "",
		"repo/.git/hooks/main.tf":                 "",
		"web/node_modules/pkg/main.tf":            "",
		"cache/.terraform/modules/remote/main.tf": "",
	})

	cfg := config.DefaultConfig()
	cfg.SearchPaths = []string{root, root} // a repeated search path must not duplicate projects

	projects, err := DiscoverProjects(cfg)
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]Project{}
	for _, p := range projects {
		rel, _ := filepath.Rel(root, p.Path)
		got[rel] = p
	}
	wantPaths := []string{"live/app", "live/network", "other/app"}
	var gotPaths []string
	for rel := range got {
		gotPaths = append(gotPaths, filepath.ToSlash(rel))
	}
	slices.Sort(gotPaths)
	if !slices.Equal(gotPaths, wantPaths) {
		t.Fatalf("discovered %v, want %v", gotPaths, wantPaths)
	}

	if !got[filepath.FromSlash("live/app")].IsInitialized {
		t.Error("live/app has .terraform/, want IsInitialized")
	}
	if got[filepath.FromSlash("live/network")].IsInitialized {
		t.Error("live/network has no .terraform/, want not initialized")
	}
	if name := got[filepath.FromSlash("live/app")].Name; name != "live/app" {
		t.Errorf("duplicate basename: Name = %q, want live/app", name)
	}
	if name := got[filepath.FromSlash("live/network")].Name; name != "network" {
		t.Errorf("unique basename: Name = %q, want network", name)
	}
}

func TestDiscoverProjectsMissingSearchPath(t *testing.T) {
	cfg := config.Config{SearchPaths: []string{filepath.Join(t.TempDir(), "missing")}}
	projects, err := DiscoverProjects(cfg)
	if err != nil || len(projects) != 0 {
		t.Errorf("DiscoverProjects() = %v, %v; want no projects and no error", projects, err)
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	tests := []struct{ in, want string }{
		{"~/Projects", filepath.Join(home, "Projects")},
		{"~", home},
		{"/abs/path", "/abs/path"},
		{"rel/~path", "rel/~path"},
	}
	for _, tt := range tests {
		got, err := ExpandPath(tt.in)
		if err != nil || filepath.Clean(got) != filepath.Clean(tt.want) {
			t.Errorf("ExpandPath(%q) = %q, %v; want %q", tt.in, got, err, tt.want)
		}
	}
}

func TestIsTerraformProject(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"tf/main.tf":          "",
		"json/main.tf.json":   "{}",
		"none/README.md":      "",
		"notjson/x.json":      "{}", // plain JSON is not Terraform
		"jsondir/a.tf.json/x": "",   // a directory named *.tf.json is not a file
		"dironly/sub.tf/x.md": "",   // a directory named *.tf is not a .tf file
	})
	tests := []struct {
		dir  string
		want bool
	}{
		{"tf", true},
		{"json", true},
		{"notjson", false},
		{"jsondir", false},
		{"none", false},
		{"dironly", false},
	}
	for _, tt := range tests {
		got, err := IsTerraformProject(filepath.Join(root, tt.dir))
		if err != nil || got != tt.want {
			t.Errorf("IsTerraformProject(%s) = %v, %v; want %v", tt.dir, got, err, tt.want)
		}
	}
}
