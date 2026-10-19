package terraform

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestBuildPlanArgs(t *testing.T) {
	// The table: each row is one named scenario with its inputs and expected output.
	tests := []struct {
		name     string
		varFile  VarFile
		planFile string
		want     []string
	}{
		{
			name:     "with var file",
			varFile:  VarFile{FullPath: "/home/user/infra/dev.tfvars"},
			planFile: "/cache/infra-dev.tfplan",
			want:     []string{"plan", "-input=false", "-var-file=/home/user/infra/dev.tfvars", "-out=/cache/infra-dev.tfplan"},
		},
		{
			name:     "without var file",
			varFile:  VarFile{}, // empty FullPath -> no -var-file flag
			planFile: "/cache/infra-default.tfplan",
			want:     []string{"plan", "-input=false", "-out=/cache/infra-default.tfplan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPlanArgs(tt.varFile, tt.planFile)
			if !slices.Equal(got, tt.want) {
				t.Errorf("buildPlanArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanFilePathIsolation(t *testing.T) {
	cache := "/cache"
	devA := planFilePath(cache, "/work/a/infra", "dev")

	tests := []struct {
		name        string
		projectPath string
		envName     string
		wantSame    bool
	}{
		{name: "same project and env", projectPath: "/work/a/infra", envName: "dev", wantSame: true},
		{name: "unclean path to same project", projectPath: "/work/a/../a/infra/", envName: "dev", wantSame: true},
		{name: "other env of same project", projectPath: "/work/a/infra", envName: "prod", wantSame: false},
		{name: "same basename, other project", projectPath: "/work/b/infra", envName: "dev", wantSame: false},
		{name: "no env", projectPath: "/work/a/infra", envName: "", wantSame: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := planFilePath(cache, tt.projectPath, tt.envName)
			if (got == devA) != tt.wantSame {
				t.Errorf("planFilePath(%q, %q) = %q, compared to %q: same=%v, want %v",
					tt.projectPath, tt.envName, got, devA, got == devA, tt.wantSame)
			}
			if filepath.Dir(got) != cache {
				t.Errorf("planFilePath() = %q, want it directly under %q", got, cache)
			}
		})
	}
}

func TestPlanFilePathEnvCannotEscapeCacheDir(t *testing.T) {
	got := planFilePath("/cache", "/work/infra", "../../etc/dev")
	if filepath.Dir(got) != "/cache" {
		t.Errorf("planFilePath() = %q, want it directly under /cache", got)
	}
}

func TestPlanFilePathRelativeProject(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	rel := planFilePath("/cache", ".", "dev")
	abs := planFilePath("/cache", wd, "dev")
	if rel != abs {
		t.Errorf("relative path gave %q, absolute gave %q; want the same plan file", rel, abs)
	}
}

func TestPlanOlderThanInit(t *testing.T) {
	project := t.TempDir()
	initTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	if PlanOlderThanInit(project, initTime.Add(-time.Hour)) {
		t.Error("PlanOlderThanInit() = true on a project never initialised, want false")
	}

	stateFile := filepath.Join(project, ".terraform", "terraform.tfstate")
	if err := os.MkdirAll(filepath.Dir(stateFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stateFile, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(stateFile, initTime, initTime); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		planTime time.Time
		want     bool
	}{
		{name: "plan before init", planTime: initTime.Add(-time.Minute), want: true},
		{name: "plan after init", planTime: initTime.Add(time.Minute), want: false},
		{name: "plan at init", planTime: initTime, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlanOlderThanInit(project, tt.planTime); got != tt.want {
				t.Errorf("PlanOlderThanInit() = %v, want %v", got, tt.want)
			}
		})
	}
}
